package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"imagu/auth"
	"imagu/db"
	"imagu/db/repo"
	"imagu/routes"
	"imagu/util"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.InitDB("./database.db")
	defer db.CloseDB()

	err = repo.InitUserRepo()
	if err != nil {
		log.Fatal("Error init user repo")
	}
	err = repo.InitImageRepo()
	if err != nil {
		log.Fatal("Error init image repo")
	}

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20

	r.StaticFile("/", "./static/html/index.html")
	r.StaticFile("/login", "./static/html/login.html")
	r.StaticFile("/register", "./static/html/register.html")
	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")

	r.Use(auth.AuthMiddleware())

	r.POST("/api/v1/login", routes.LoginUser)
	r.POST("/api/v1/user/checkPermission", routes.CheckPermission)
	r.POST("/api/v1/register", auth.AuthPermission("register"), routes.RegisterUser)

	r.POST("/api/v1/image/uploadImage", auth.AuthPermission("uploadImage"), routes.UploadImage)
	r.GET("/image/get/:id", auth.AuthPermission("viewImage"), routes.GetImage)
	r.GET("/image/view/:id", auth.AuthPermission("viewImage"), routes.ViewImage)

	r.GET("/api/v1/stats", routes.GetStats)

	r.Run(":8080")
	util.CloseLog()
}
