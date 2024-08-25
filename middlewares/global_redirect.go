package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"imagu/db/repo"
	"net/http"
)

// GlobalRedirect redirects everything to the admin_register until an admin was created
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
