package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"imagu/internal/db"
	repo2 "imagu/internal/db/repo"
	"imagu/internal/jobs"
	"imagu/internal/middlewares/auth"
	routes2 "imagu/internal/routes"
	"imagu/internal/stats"
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

	err := repo2.InitUserRepo()
	if err != nil {
		logrus.Fatal("Error init user repo: ", err)
	}
	err = repo2.InitImageRepo()
	if err != nil {
		logrus.Fatal("Error init image repo: ", err)
	}
	err = repo2.InitSettingsRepo()
	if err != nil {
		logrus.Fatal("Error init settings repo: ", err)
	}
}
func main() {
	requiredEnvs()
	createDirs()
	file, err := os.OpenFile("./logs/"+time.Now().Format(time.DateTime)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(io.MultiWriter(file, os.Stdout))
	logrus.SetLevel(logrus.InfoLevel)
	initDB()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(LoggerWithLogrus(), gin.Recovery())
	r.MaxMultipartMemory = 8 << 20
	r.Use(auth.AuthMiddleware())

	routes2.InitSiteRoutes(r)
	routes2.InitUserRoutes(r)
	routes2.InitStatsRoutes(r)
	routes2.InitImageRoutes(r)
	routes2.InitAdminRoutes(r)

	jobHandler := jobs.JobHandler{}
	jobHandler.AddJob(jobs.DeletionJob)
	go jobHandler.Run()
	logrus.Info("Starting Gin server")
	err = r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		logrus.Fatal("Could not Run gin: ", err)
	}
	db.CloseDB()
	stats.CloseLog()
}
