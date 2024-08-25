package routes

import (
	"github.com/gin-gonic/gin"
	"imagu/db/repo"
	"imagu/util"
)

type StatsResponse struct {
	Images      int              `json:"images"`
	SubImages   int              `json:"sub_images"`
	TotalImages int              `json:"total_images"`
	ImageSize   string           `json:"image_size"`
	AccessStats util.AccessStats `json:"access_stats"`
}

func GetStats(c *gin.Context) {
	images, _ := repo.GetAllImages()
	imageCount := len(images)
	subImageCount := 0
	var imagesSize int64 = 0
	for _, value := range images {
		imagesSize += value.Size
		subImageCount += value.SubImages
	}
	stats := StatsResponse{
		Images:      imageCount,
		SubImages:   subImageCount,
		TotalImages: imageCount + subImageCount,
		ImageSize:   util.FormatBytesToString(imagesSize),
		AccessStats: *util.CurrentAccessStats,
	}
	c.JSON(200, stats)
}
