package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"html/template"
	"imagu/internal/converter"
	"imagu/internal/db/repo"
	"imagu/internal/middlewares/auth"
	"imagu/internal/stats"
	util2 "imagu/internal/util"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var uploadsDir = "./data/uploads"

// InitImageRoutes sets up the routes for image-related endpoints.
func InitImageRoutes(r *gin.Engine) {
	r.POST("/api/v1/image/uploadImage", auth.AuthPermission("uploadImage", false), uploadImage)
	r.GET("/image/get/:id", auth.AuthPermission("viewImage", false), getImage)
	r.GET("/image/view/:id", auth.AuthPermission("viewImage", false), viewImage)
}

// isImageTypeLossy determines if the provided file extension indicates a lossy image type.
func isImageTypeLossy(ext string) bool {
	lossyExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".webp": true,
	}
	return lossyExtensions[ext]
}

// getBaseImageFileName retrieves the name of the base image file from the specified directory.
func getBaseImageFileName(dir string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "base") {
			return file.Name(), nil
		}
	}
	return "", fmt.Errorf("base file not found")
}

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
	allowedExtensions := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".webp": true,
		".qoi":  true,
	}
	fileExt := strings.ToLower(filepath.Ext(filename))
	return allowedExtensions[fileExt]
}

// uploadImage handles image uploads, validates the image, and performs conversions if necessary.
func uploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		logrus.Warn("No file received")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file received"})
		return
	}
	if !isValidImageType(file) || !isValidImageExtension(file.Filename) {
		logrus.Warn("Invalid file type. Supported formats: png, jpg, jpeg, webp, qoi")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Supported formats: png, jpg, jpeg, webp, qoi"})
		return
	}

	user := util2.GetUserFromContext(c)
	if user == nil {
		logrus.Error("User not found")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	id := uuid.New()
	directory := filepath.Join(uploadsDir, id.String())
	if err = os.MkdirAll(directory, 0755); err != nil {
		logrus.Errorf("Failed to create directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	dst := filepath.Join(directory, file.Filename)
	if err = c.SaveUploadedFile(file, dst); err != nil {
		logrus.Errorf("Failed to upload file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	var baseFile string
	if isImageTypeLossy(filepath.Ext(dst)) {
		baseFile = "base.webp"
		err = converter.ConvertImage(dst, filepath.Join(directory, baseFile), "", "", "", "")
	} else {
		baseFile = "base.qoi"
		err = converter.ConvertImage(dst, filepath.Join(directory, baseFile), "", "", "", "")
	}

	if err != nil {
		logrus.Errorf("Failed to convert image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert image"})
		return
	}

	stat, err := os.Stat(filepath.Join(directory, baseFile))
	if err != nil {
		logrus.Errorf("Failed to retrieve file data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file data"})
		return
	}

	if err = repo.CreateImage(file.Filename, user.ID, id.String(), stat.Size()); err != nil {
		logrus.Errorf("Failed to create image record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image record"})
		return
	}

	if err = os.Remove(dst); err != nil {
		logrus.Warnf("Could not delete uploaded file: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"url":     fmt.Sprintf("/image/view/%s", id.String()),
	})
}

// getImage retrieves the image based on the provided ID, applying transformations if needed.
func getImage(c *gin.Context) {
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

	filePath := fmt.Sprintf("%s/%s/alt_", uploadsDir, uuid)
	wasChanged := false

	if quality != "" {
		filePath += fmt.Sprintf("q:%s", quality)
		wasChanged = true
	}
	if resizeQuery != "" && crop == "" {
		filePath += fmt.Sprintf("r:%s", resizeQuery)
		wasChanged = true
	}
	if crop != "" {
		filePath += fmt.Sprintf("c:%s", crop)
		wasChanged = true
	}
	if blur != "" {
		filePath += fmt.Sprintf("b:%s", blur)
		wasChanged = true
	}

	if !wasChanged {
		filePath = fmt.Sprintf("%s/%s/alt", uploadsDir, uuid)
	}
	filePath += fmt.Sprintf(".%s", ext)

	if !util2.FileExists(filePath) {
		dir := filepath.Join(uploadsDir, uuid)
		baseFile, err := getBaseImageFileName(dir)
		if err != nil {
			logrus.Errorf("Failed to get base image: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base image"})
			return
		}
		err = converter.ConvertImage(filepath.Join(dir, baseFile), filePath, resizeQuery, quality, crop, blur)
		if err != nil {
			logrus.Errorf("Failed to create image: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image"})
			return
		}
		stat, _ := os.Stat(filePath)
		if err := repo.UpdateSizeAndCount(uuid, stat.Size(), 1); err != nil {
			logrus.Warnf("Could not update image size: %v", err)
		}
	}

	stats.LogAccess(uuid)
	c.File(filePath)
}

type viewImagePage struct {
	Title string
}

// viewImage displays the image view page for the given image ID.
func viewImage(c *gin.Context) {
	idParam := c.Param("id")
	image, err := repo.GetImageFromUUID(idParam)
	if err != nil {
		logrus.Warnf("Image with ID %s not found", idParam)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	viewpageData := viewImagePage{
		Title: image.Name,
	}
	tmplPath := path.Join("static", "html", "image", "imageView.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		logrus.Errorf("Could not parse template %s: %v", tmplPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse template"})
		return
	}

	if err := tmpl.Execute(c.Writer, viewpageData); err != nil {
		logrus.Errorf("Could not execute template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not execute template"})
		return
	}
}
