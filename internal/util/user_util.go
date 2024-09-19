package util

import (
	"github.com/gin-gonic/gin"
	"imagu/internal/db/model"
)

func GetUserFromContext(c *gin.Context) *model.User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	u, ok := user.(*model.User)
	if !ok {
		return nil
	}
	return u
}
