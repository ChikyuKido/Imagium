package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"html/template"
	"imagu/converter"
	"imagu/db/repo"
	"imagu/util"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var uploadsDir = "./data/uploads"
var baseFile = "base.qoi"

// isValidImageType checks if the uploaded file has an allowed MIME type.
func isValidImageType(fileHeader *multipart.FileHeader) bool {
	allowedMIMETypes := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
		"image/webp": true,
		"image/qoi":  true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	return allowedMIMETypes[contentType]
}

// isValidImageExtension checks if the file extension is among the allowed ones.
func isValidImageExtension(filename string) bool {
	allowedExtensions := []string{".png", ".jpg", ".jpeg", ".webp", ".qoi"}
	if len(strings.Split(filename, ".")) != 2 {
		return false
	}
	fileExt := strings.ToLower(filename[strings.LastIndex(filename, "."):])

	for _, ext := range allowedExtensions {
		if fileExt == ext {
			return true
		}
	}
	return false
}

// UploadImage handles image uploads, validates them, and performs conversions.
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file received"})
		return
	}
	if !isValidImageType(file) || !isValidImageExtension(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Supported formats: png, jpg, jpeg, webp, qoi"})
		return
	}
	user := util.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	id := uuid.New()
	err = os.MkdirAll(filepath.Join(uploadsDir, id.String()), 0755)
	if err != nil {
		logrus.Errorf("Failed to create directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}
	dst := filepath.Join(uploadsDir, id.String(), file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		logrus.Errorf("Failed to upload file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}
	err = converter.ConvertImage(dst, filepath.Join(filepath.Dir(dst), baseFile), "", "", "", "")
	if err != nil {
		logrus.Errorf("Failed to convert image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert image"})
		return
	}
	stat, err := os.Stat(filepath.Join(filepath.Dir(dst), baseFile))
	if err != nil {
		logrus.Errorf("Failed to retrieve file data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file data"})
		return
	}
	err = repo.CreateImage(file.Filename, user.ID, id.String(), stat.Size())
	if err != nil {
		logrus.Errorf("Failed to create image record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image record"})
		return
	}
	err = os.Remove(dst)
	if err != nil {
		logrus.Warnf("Could not delete uploaded file: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"url":     fmt.Sprintf("/image/view/%s", id.String()),
	})
}

// GetImage handles requests to retrieve images, applying transformations if needed.
func GetImage(c *gin.Context) {
	idParam := c.Param("id")
	if !isValidImageExtension(idParam) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Invalid file extension"})
		return
	}
	resizeQuery := c.Query("resize")
	quality := c.Query("quality")
	blur := c.Query("blur")
	crop := strings.Replace(c.Query("crop"), " ", "+", -1)

	uuid := strings.Split(idParam, ".")[0]
	ext := filepath.Ext(idParam)[1:]

	file := fmt.Sprintf("%s/%s/base_", uploadsDir, uuid)
	wasChanged := false
	if quality != "" {
		file += fmt.Sprintf("q:%s", quality)
		wasChanged = true
	}
	if resizeQuery != "" && crop == "" {
		file += fmt.Sprintf("r:%s", resizeQuery)
		wasChanged = true
	}
	if crop != "" {
		file += fmt.Sprintf("c:%s", crop)
		wasChanged = true
	}
	if blur != "" {
		file += fmt.Sprintf("b:%s", blur)
		wasChanged = true
	}
	if !wasChanged {
		file = fmt.Sprintf("%s/%s/base", uploadsDir, uuid)
	}
	file += fmt.Sprintf(".%s", ext)
	if !util.FileExists(file) {
		err := converter.ConvertImage(filepath.Join(uploadsDir, uuid, baseFile), file, resizeQuery, quality, crop, blur)
		if err != nil {
			logrus.Errorf("Failed to create image: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image"})
			return
		}
		stat, _ := os.Stat(file)
		err = repo.UpdateSizeAndCount(uuid, stat.Size(), 1)
		if err != nil {
			logrus.Warnf("Could not update image size: %v", err)
		}
	}
	util.LogAccess(uuid)
	c.File(file)
}

type ViewImagePage struct {
	Title string
}

// ViewImage handles the display of the image view page.
func ViewImage(c *gin.Context) {
	idParam := c.Param("id")
	image, err := repo.GetImageFromUUID(idParam)
	if err != nil {
		c.Redirect(http.StatusPermanentRedirect, "/")
		return
	}

	viewpageData := ViewImagePage{
		Title: image.Name,
	}
	tmpl, err := template.ParseFiles(path.Join("static", "html", "image", "imageView.html"))
	if err != nil {
		logrus.Errorf("Could not parse template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse template"})
		return
	}

	err = tmpl.Execute(c.Writer, viewpageData)
	if err != nil {
		logrus.Errorf("Could not execute template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not execute template"})
		return
	}
}
