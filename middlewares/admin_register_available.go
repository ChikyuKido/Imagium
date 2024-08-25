package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"imagu/db/repo"
	"net/http"
)

// AdminRegisterAvailable Checks if the admin registration is still available. Only goes to the next stage if no admin was created yet
func AdminRegisterAvailable(redirect bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		present, err := repo.GetAdminUser()
		if err != nil {
			if redirect {
				c.Redirect(http.StatusPermanentRedirect, "/login")
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check admin user"})
			}
			c.Abort()
			logrus.Error("Could not check if the admin user exists: ", err)
		}
		if present {
			if redirect {
				c.Redirect(http.StatusPermanentRedirect, "/login")
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin user already exists"})
			}
			c.Abort()
		} else {
			c.Next()
		}
	}
}
