package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"imagu/db/model"
	"imagu/db/repo"
	"imagu/util"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || tokenString == "guest" {
			guest, err := repo.GetUserByName("guest")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No guest user found."})
				c.Abort()
				return
			}
			c.Set("user", guest)
			c.Next()
			return
		}

		token, err := util.GetToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		userName := claims["name"].(string)
		user, err := repo.GetUserByName(userName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
func AuthPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		if repo.HasPermission(u, permission) {
			c.Next()
		} else {
			c.Redirect(http.StatusPermanentRedirect, "/login")
		}
	}
}
func GlobalRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/admin_register" {
			c.Next()
		}
		present, err := repo.GetAdminUser()
		if err != nil {
			logrus.Error("Could not check if the admin user exists: ", err)
			c.Redirect(http.StatusPermanentRedirect, "/admin_register")
			c.Abort()
		}
		if !present {
			c.Redirect(http.StatusPermanentRedirect, "/admin_register")
			c.Abort()
		}
	}
}
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
