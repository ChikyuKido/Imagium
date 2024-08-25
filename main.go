package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"imagu/db"
	"imagu/db/repo"
	"imagu/middlewares"
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

func existEnv(key string) {
	myVar := os.Getenv(key)
	if myVar == "" {
		panic("env variable MY_VAR not set: " + key)
	}
}
func requiredEnvs() {
	existEnv("DOMAIN")
	existEnv("PORT")
}
func createDirs() {
	err := os.MkdirAll("./logs", 0755)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll("./data", 0755)
	if err != nil {
		panic(err)
	}
}
func initDB() {
	db.InitDB("./data/database.db")
	logrus.Info("Initialized database")
	defer db.CloseDB()

	err := repo.InitUserRepo()
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
}
func initLogrus() {
	file, err := os.OpenFile("./logs/"+time.Now().Format(time.DateTime)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(io.MultiWriter(file, os.Stdout))
	logrus.SetLevel(logrus.InfoLevel)

}
func main() {
	requiredEnvs()
	createDirs()
	initLogrus()
	initDB()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(LoggerWithLogrus(), gin.Recovery())
	r.MaxMultipartMemory = 8 << 20
	r.Use(middlewares.AuthMiddleware())

	routes.InitSiteRoutes(r)
	routes.InitUserRoutes(r)
	routes.InitStatsRoutes(r)
	routes.InitImageRoutes(r)
	routes.InitAdminRoutes(r)

	logrus.Info("Starting Gin server")
	err := r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		logrus.Fatal("Could not Run gin: ", err)
	}
	db.CloseDB()
	util.CloseLog()
}
