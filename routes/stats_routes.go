package routes

import (
	"github.com/gin-gonic/gin"
	"imagu/util"
)

func GetStats(c *gin.Context) {
	c.JSON(200, util.CurrentAccessStats)
}
