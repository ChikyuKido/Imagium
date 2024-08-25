package routes

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"imagu/middlewares"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"sync"
)

var (
	cache      = make(map[string][]byte)
	cacheMutex = &sync.Mutex{}
)

func servePage(path string, diskPath string, r *gin.RouterGroup) {
	r.GET(path, func(c *gin.Context) {
		content := getCachedContent(path, diskPath)
		contentType := mime.TypeByExtension(filepath.Ext(diskPath))
		c.Header("Content-Encoding", "gzip")
		c.Data(200, contentType, content)
	})
}
func servePageWith(path string, diskPath string, r *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	r.GET(path, append(handlers, func(c *gin.Context) {
		content := getCachedContent(path, diskPath)
		contentType := mime.TypeByExtension(filepath.Ext(diskPath))
		c.Header("Content-Encoding", "gzip")
		c.Data(200, contentType, content)
	})...)
}
func getCachedContent(path string, filepath string) []byte {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if content, found := cache[path]; found {
		return content
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil
	}

	var compressedContent bytes.Buffer
	writer, _ := gzip.NewWriterLevel(&compressedContent, gzip.BestCompression)
	_, err = writer.Write(content)
	if err != nil {
		return nil
	}
	writer.Close()

	compressedData := compressedContent.Bytes()
	cache[path] = compressedData
	return compressedData
}

func serveDirectory(rootPath string, baseDir string, r *gin.RouterGroup) {
	filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(baseDir, path)
			urlPath := rootPath + relativePath
			servePage(urlPath, path, r)
		}
		return nil
	})
}

func InitSiteRoutes(r *gin.Engine) {
	sitesGroup := r.Group("/")
	redirectGroup := sitesGroup.Group("/")
	redirectGroup.Use(middlewares.GlobalRedirect())

	servePageWith("/", "./static/html/index.html", redirectGroup, middlewares.AuthPermission("uploadImage", true))
	servePageWith("/register", "./static/html/register.html", redirectGroup, middlewares.AuthPermission("register", true))
	servePage("/login", "./static/html/login.html", redirectGroup)
	servePageWith("/admin/register", "./static/html/admin/register.html", redirectGroup, middlewares.AdminRegisterAvailable(true))
	servePageWith("/admin/dashboard", "./static/html/admin/dashboard.html", redirectGroup, middlewares.AuthPermission("admin", true))

	serveDirectory("/js/", "./static/js", sitesGroup)
	serveDirectory("/css/", "./static/css", sitesGroup)
	serveDirectory("/imgs/", "./static/imgs", sitesGroup)
}
