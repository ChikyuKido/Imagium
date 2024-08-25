package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"imagu/db/repo"
	"net/http"
)

func AdminRegister(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create admin user"})
		return
	}
	user, err := repo.GetUserByName(request.Username)
	if err != nil {
		// could not create user because if it cant find the user it was not created.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not admin create user"})
		logrus.Error("Admin user not found: ", err)
		return
	}

	err = repo.AddPermission(user, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add role admin"})
		logrus.Error("Could not add role admin: ", err)
		return
	}
	err = repo.AddPermission(user, "register")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add role register"})
		logrus.Error("Could not add role register: ", err)
		return
	}
	err = repo.UpdateAdminUser(true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could update admin user settings value. This means another admin user can be created"})
		logrus.Error("Could update admin user settings value. This means another admin user can be created", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully register a admin user"})
}
