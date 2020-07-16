package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"github.com/shinshin86/go-todo/db"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	db.Init()

	router.GET("/", func(ctx *gin.Context) {
		todos := db.GetAll()
		fmt.Println(todos)
		ctx.HTML(200, "index.html", gin.H{
			"todos": todos,
		})
	})

	router.POST("/new", func(ctx *gin.Context) {
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		tag := ctx.PostForm("tag")
		db.Insert(text, status, tag)
		ctx.Redirect(302, "/")
	})

	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}

		todo := db.GetOne(id)
		ctx.HTML(200, "detail.html", gin.H{"todo": todo})
	})

	router.GET("/tag/:tag", func(ctx *gin.Context) {
		tagname := ctx.Param("tag")

		tag := db.GetTagItem(tagname)
		fmt.Println(tag)
		ctx.HTML(200, "tag.html", gin.H{"tag": tag})
	})

	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		tag := ctx.PostForm("tag")
		db.Update(id, text, status, tag)
		ctx.Redirect(302, "/")
	})

	router.POST("/delete", func(ctx *gin.Context) {
		n := ctx.PostForm("delete-id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		db.Delete(id)
		ctx.Redirect(302, "/")
	})

	router.Run()
}
