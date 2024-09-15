package routes

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"imagu/db/model"
	"imagu/db/repo"
	"imagu/middlewares"
	"imagu/util"
	"net/http"
	"os"
	"strconv"
)

func InitUserRoutes(r *gin.Engine) {
	r.POST("/api/v1/user/login", loginUser)
	r.POST("/api/v1/user/register", middlewares.AuthPermission("register", false), registerUser)
	r.GET("/api/v1/user/library", middlewares.AuthPermission("viewLibrary", false), library)
}

func registerUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if repo.DoesUserByNameExists(request.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}
	err := repo.CreateUser(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
func loginUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := repo.GetUserByName(request.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !util.CheckPasswordHash(request.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
		return
	}

	token, err := util.GenerateJWT(request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.SetCookie("jwt", token, 60*60*24*30, "/", os.Getenv("DOMAIN"), false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in"})
}
func library(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		c.Abort()
		return
	}
	u, ok := user.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		c.Abort()
		return
	}
	siteQuery, _ := c.GetQuery("site")
	site := 0
	if siteQuery != "" {
		s, err := strconv.Atoi(siteQuery)
		if err != nil {
			site = 0
		} else {
			site = s
		}
	}
	images, err := repo.GetAllImagesByUser(u, 20*site, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user"})
		c.Abort()
		return
	}
	count, err := repo.GetImageCountByUser(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user"})
		c.Abort()
		return
	}

	var response struct {
		Images []model.ImageModel
		Pages  int
	}
	response.Images = images
	response.Pages = int(count / 20)

	c.JSON(http.StatusOK, response)
}
