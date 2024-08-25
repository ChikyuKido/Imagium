package routes

import (
	"github.com/gin-gonic/gin"
	"imagu/middlewares"
)

func InitSiteRoutes(r *gin.Engine) {
	sitesGroup := r.Group("/")
	sitesGroup.Use(middlewares.GlobalRedirect())
	// only available if the user has upload permission
	sitesGroup.GET("/", middlewares.AuthPermission("uploadImage", true), func(c *gin.Context) {
		c.File("./static/html/index.html")
	})
	// always accessible
	sitesGroup.StaticFile("/login", "./static/html/login.html")
	// only available if the user has register permissions
	sitesGroup.GET("/register", middlewares.AuthPermission("register", true), func(c *gin.Context) {
		c.File("./static/html/register.html")
	})
	// only is available if no admin user was created
	sitesGroup.GET("/admin_register", middlewares.AdminRegisterAvailable(true), func(c *gin.Context) {
		c.File("./static/html/admin_register.html")
	})
	sitesGroup.GET("/admin/dashboard", middlewares.AuthPermission("admin", true), func(c *gin.Context) {
		c.File("./static/html/admin/dashboard.html")
	})

	// static resources
	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")
}
