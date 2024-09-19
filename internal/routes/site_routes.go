package routes

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	middlewares2 "imagu/internal/middlewares"
	"imagu/internal/middlewares/auth"
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
	redirectGroup.Use(middlewares2.GlobalRedirect())

	servePageWith("/", "./static/html/index.html", redirectGroup, auth.AuthPermission("uploadImage", true))
	servePageWith("/register", "./static/html/register.html", redirectGroup, auth.AuthPermission("register", true))
	servePageWith("/library", "./static/html/library.html", redirectGroup, auth.AuthPermission("viewLibrary", true))
	servePage("/login", "./static/html/login.html", redirectGroup)
	servePageWith("/admin/register", "./static/html/admin/register.html", redirectGroup, middlewares2.AdminRegisterAvailable(true))
	servePageWith("/admin/dashboard", "./static/html/admin/dashboard.html", redirectGroup, auth.AuthPermission("admin", true))

	serveDirectory("/js/", "./static/js", sitesGroup)
	serveDirectory("/css/", "./static/css", sitesGroup)
	serveDirectory("/imgs/", "./static/imgs", sitesGroup)
}
