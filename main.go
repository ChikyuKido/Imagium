package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"imagu/auth"
	"imagu/db"
	"imagu/db/repo"
	"imagu/routes"
	"imagu/util"
	"io"
	"os"
	"time"
)

func LoggerWithLogrus() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logrus.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  c.Writer.Status(),
			"latency": time.Since(start),
		}).Debug("Request handled")
	}
}

func main() {
	err := os.MkdirAll("./logs", 0755)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll("./data", 0755)
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile("./logs/"+time.Now().Format(time.DateTime)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(io.MultiWriter(file, os.Stdout))
	logrus.SetLevel(logrus.InfoLevel)

	db.InitDB("./data/database.db")
	logrus.Info("Initialized database")
	defer db.CloseDB()

	err = repo.InitUserRepo()
	if err != nil {
		logrus.Fatal("Error init user repo: ", err)
	}
	err = repo.InitImageRepo()
	if err != nil {
		logrus.Fatal("Error init image repo: ", err)
	}
	err = repo.InitSettingsRepo()
	if err != nil {
		logrus.Fatal("Error init settings repo: ", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(LoggerWithLogrus(), gin.Recovery())
	r.MaxMultipartMemory = 8 << 20

	r.Use(auth.AuthMiddleware())

	sitesGroup := r.Group("/")
	sitesGroup.Use(auth.GlobalRedirect())

	// pages
	sitesGroup.GET("/", auth.AuthPermission("uploadImage"), func(c *gin.Context) {
		c.File("./static/html/index.html")
	})
	sitesGroup.StaticFile("/login", "./static/html/login.html")
	sitesGroup.GET("/register", auth.AuthPermission("register"), func(c *gin.Context) {
		c.File("./static/html/register.html")
	})
	sitesGroup.GET("/admin_register", auth.AdminRegisterAvailable(true), func(c *gin.Context) {
		c.File("./static/html/admin_register.html")
	})

	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")

	// api endpoints
	r.POST("/api/v1/login", routes.LoginUser)
	r.POST("/api/v1/user/checkPermission", routes.CheckPermission)
	r.POST("/api/v1/register", auth.AuthPermission("register"), routes.RegisterUser)

	r.POST("/api/v1/image/uploadImage", auth.AuthPermission("uploadImage"), routes.UploadImage)
	r.GET("/image/get/:id", auth.AuthPermission("viewImage"), routes.GetImage)
	r.GET("/image/view/:id", auth.AuthPermission("viewImage"), routes.ViewImage)

	r.GET("/api/v1/stats", auth.AuthPermission("viewStats"), routes.GetStats)

	r.POST("/api/v1/admin/register", auth.AdminRegisterAvailable(false), routes.AdminRegister)

	logrus.Info("Starting Gin server")
	err = r.Run(":8080")
	if err != nil {
		logrus.Fatal("Could not Run gin: ", err)
	}
	db.CloseDB()
	util.CloseLog()
}
