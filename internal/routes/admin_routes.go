package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"imagu/internal/db/model"
	repo2 "imagu/internal/db/repo"
	"imagu/internal/middlewares"
	"imagu/internal/middlewares/auth"
	"net/http"
	"strings"
)

func InitAdminRoutes(r *gin.Engine) {
	r.POST("/api/v1/admin/register", middlewares.AdminRegisterAvailable(false), adminRegister)
	r.GET("/api/v1/admin/users", auth.AuthPermission("admin", false), getAllUsers)
	r.PUT("/api/v1/admin/users/changeRole/:userId", auth.AuthPermission("admin", false), changeRole)
	r.GET("/api/v1/admin/settings", auth.AuthPermission("admin", false), settings)
	r.PUT("/api/v1/admin/updateSettings", auth.AuthPermission("admin", false), updateSettings)
}

func updateSettings(c *gin.Context) {
	var settings model.SettingsModel
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo2.UpdateSettings(settings)

}

func settings(c *gin.Context) {
	settings, err := repo2.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func adminRegister(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if repo2.DoesUserByNameExists(request.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}
	err := repo2.CreateUser(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create admin user"})
		return
	}
	user, err := repo2.GetUserByName(request.Username)
	if err != nil {
		// could not create user because if it cant find the user it was not created.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not admin create user"})
		logrus.Error("Admin user not found: ", err)
		return
	}

	user.Roles = strings.Join(append(strings.Split(user.Roles, ","), "register", "admin"), ",")
	err = repo2.UpdateRole(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add roles to admin"})
		logrus.Error("Could not add roles to admin: ", err)
		return
	}
	err = repo2.UpdateAdminUser(true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could update admin user settings value. This means another admin user can be created"})
		logrus.Error("Could update admin user settings value. This means another admin user can be created", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully register a admin user"})
}
func getAllUsers(c *gin.Context) {
	users, err := repo2.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
func changeRole(c *gin.Context) {
	userId := c.Param("userId")
	var request struct {
		Roles string `json:"roles" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := repo2.GetUserById(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not find user with id: " + userId})
		return
	}

	user.Roles = request.Roles
	err = repo2.UpdateRole(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not update roles for user: " + userId})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "changed user roles"})
}
