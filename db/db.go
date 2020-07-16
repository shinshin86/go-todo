package db

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/jinzhu/gorm"
)

type Todo struct {
	gorm.Model
	Text   string
	Status string
	TagID  uint
	Tag    Tag
}

type Tag struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Color string
	Todos []Todo
}

func dbOpen() *gorm.DB {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("Not open database error")
	}

	return db
}

func getRandomColorCode() (string, error) {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func Init() {
	db := dbOpen()
	db.AutoMigrate(&Todo{})
	db.AutoMigrate(&Tag{})
	defer db.Close()
}

func Insert(text string, status string, tag string) {
	db := dbOpen()

	var existsTag Tag
	db.Where("Name = ?", tag).First(&existsTag)

	if existsTag.Name != "" {
		todo := &Todo{Text: text, Status: status, Tag: existsTag}
		db.Create(&todo)
	} else {
		color, err := getRandomColorCode()
		if err != nil {
			panic("Create color code error")
		}

		todo := &Todo{
			Text:   text,
			Status: status,
			Tag: Tag{
				Name:  tag,
				Color: "#" + color,
			}}
		db.Create(&todo)
	}

	defer db.Close()
}

func Update(id int, text string, status string, tag string) {
	db := dbOpen()

	var todo Todo
	var tag1 Tag

	db.First(&todo, id).Model(&todo).Related(&tag1)

	todo.Text = text
	todo.Status = status
	tag1.Name = tag

	db.Save(&todo)
	db.Save(&tag1)
	db.Close()
}

func Delete(id int) {
	db := dbOpen()

	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)

	db.Close()
}

func GetAll() []Todo {
	db := dbOpen()

	var todos []Todo
	db.Preload("Tag").Order("created_at desc").Find(&todos)

	db.Close()
	return todos
}

func GetOne(id int) Todo {
	db := dbOpen()

	var todo Todo
	db.Preload("Tag").First(&todo, id)

	db.Close()
	return todo
}

func GetTagItem(tagname string) Tag {
	db := dbOpen()

	var tag Tag
	db.Preload("Todos").Where("Name = ?", tagname).First(&tag)

	db.Close()
	return tag
}
