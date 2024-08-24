package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"imagu/converter"
	"imagu/db/repo"
	"imagu/util"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func isValidImageType(fileHeader *multipart.FileHeader) bool {
	allowedMIMETypes := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
		"image/webp": true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	return allowedMIMETypes[contentType]
}
func isValidImageExtension(filename string) bool {
	allowedExtensions := []string{".png", ".jpg", ".jpeg", ".webp"}
	fileExt := strings.ToLower(filename[strings.LastIndex(filename, "."):])

	for _, ext := range allowedExtensions {
		if fileExt == ext {
			return true
		}
	}
	return false
}

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}
	if !isValidImageType(file) || !isValidImageExtension(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Supported are png, jpg, webp, gif"})
		return
	}
	user := util.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	id := uuid.New()
	err = os.MkdirAll("./uploads/"+id.String(), os.ModePerm)
	if err != nil {
		return
	}
	dst := "./uploads/" + id.String() + "/" + file.Filename
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}
	err = repo.CreateImage(file.Filename, user.ID, id.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	err = converter.ConvertImage(dst, filepath.Dir(dst)+"/base.png", "", "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	err = os.Remove(dst)
	if err != nil {
		fmt.Println("Could not delete downloaded file")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}
func GetImage(c *gin.Context) {
	idParam := c.Param("id")
	if !isValidImageExtension(idParam) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Extension is not allowed"})
		return
	}
	resizeQuery := c.Query("resize")
	quality := c.Query("quality")
	crop := strings.Replace(c.Query("crop"), " ", "+", -1)
	uuid := strings.Split(idParam, ".")[0]
	ext := filepath.Ext(idParam)[1:]

	file := fmt.Sprintf("./uploads/%s/base_", uuid)
	wasChanged := false
	if quality != "" {
		file += "q:" + quality
		wasChanged = true
	}
	if resizeQuery != "" && crop == "" {
		file += "r:" + resizeQuery
		wasChanged = true
	}
	if crop != "" {
		file += "c:" + crop
		wasChanged = true
	}
	if !wasChanged {
		file = fmt.Sprintf("./uploads/%s/base", uuid)
	}
	file += "." + ext
	if !util.FileExists(file) {
		err := converter.ConvertImage("./uploads/"+uuid+"/base.png", file, resizeQuery, quality, crop)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
			return
		}
	}
	c.File(file)
}
