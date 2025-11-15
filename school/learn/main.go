package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed html-css-tutorial.html example-1-personal.html example-2-business.html example-3-blog.html index.html
var embeddedFiles embed.FS

func main() {
	// تنظیم Gin
	router := gin.Default()

	// صفحه اصلی - لیست منوی انتخاب
	router.GET("/", func(c *gin.Context) {
		content, err := embeddedFiles.ReadFile("index.html")
		if err != nil {
			c.JSON(500, gin.H{"error": "خطا در بارگذاری صفحه"})
			return
		}
		c.Data(200, "text/html; charset=utf-8", content)
	})

	// صفحه آموزش
	router.GET("/tutorial", func(c *gin.Context) {
		content, err := embeddedFiles.ReadFile("html-css-tutorial.html")
		if err != nil {
			c.JSON(500, gin.H{"error": "خطا در بارگذاری آموزش"})
			return
		}
		c.Data(200, "text/html; charset=utf-8", content)
	})

	// صفحه سایت شخصی
	router.GET("/personal", func(c *gin.Context) {
		content, err := embeddedFiles.ReadFile("example-1-personal.html")
		if err != nil {
			c.JSON(500, gin.H{"error": "خطا در بارگذاری سایت شخصی"})
			return
		}
		c.Data(200, "text/html; charset=utf-8", content)
	})

	// صفحه بلاگ
	router.GET("/blog", func(c *gin.Context) {
		content, err := embeddedFiles.ReadFile("example-3-blog.html")
		if err != nil {
			c.JSON(500, gin.H{"error": "خطا در بارگذاری بلاگ"})
			return
		}
		c.Data(200, "text/html; charset=utf-8", content)
	})

	// صفحه سایت شرکتی
	router.GET("/business", func(c *gin.Context) {
		content, err := embeddedFiles.ReadFile("example-2-business.html")
		if err != nil {
			c.JSON(500, gin.H{"error": "خطا در بارگذاری سایت شرکتی"})
			return
		}
		c.Data(200, "text/html; charset=utf-8", content)
	})

	// API: دریافت IP محلی و شبکه
	router.GET("/api/network-info", func(c *gin.Context) {
		localIP := getLocalIP()
		c.JSON(200, gin.H{
			"local_ip": localIP,
			"port":     8080,
			"url":      fmt.Sprintf("http://%s:8080", localIP),
		})
	})

	// API: فهرست فایل‌ها
	router.GET("/api/files", func(c *gin.Context) {
		files := []string{
			"index.html - صفحه اصلی",
			"html-css-tutorial.html - آموزش کامل",
			"example-1-personal.html - سایت شخصی",
			"example-2-business.html - سایت شرکتی",
			"example-3-blog.html - بلاگ",
		}
		c.JSON(200, gin.H{"files": files})
	})

	// دریافت محتوای فایل‌های html
	router.GET("/api/file/:filename", func(c *gin.Context) {
		filename := c.Param("filename")

		// فقط فایل‌های معتبر
		validFiles := map[string]bool{
			"index.html":              true,
			"html-css-tutorial.html":  true,
			"example-1-personal.html": true,
			"example-2-business.html": true,
			"example-3-blog.html":     true,
		}

		if !validFiles[filename] {
			c.JSON(400, gin.H{"error": "فایل نامعتبر"})
			return
		}

		content, err := embeddedFiles.ReadFile(filename)
		if err != nil {
			c.JSON(500, gin.H{"error": "خطا در بارگذاری فایل"})
			return
		}

		c.JSON(200, gin.H{
			"filename": filename,
			"size":     len(content),
			"content":  string(content),
		})
	})

	// حالت توسعه: استفاده از فایل‌های محلی
	if os.Getenv("GIN_MODE") == "debug" {
		router.Static("/static", "./static")
	}

	// راه‌اندازی سرور
	port := "8080"
	localIP := getLocalIP()

	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║              Server Started Successfully               ║")
	fmt.Println("╠══════════════════════════════════════════════════════╣")

	fmt.Printf("║  Local:    http://localhost:%-5s                      ║\n", port)
	fmt.Printf("║  Network:  http://%-15s:%-5s              ║\n", localIP, port)

	fmt.Println("║                                                      ║")
	fmt.Println("║  Available Pages:                                    ║")
	fmt.Println("║  /              - Home Page                          ║")
	fmt.Println("║  /tutorial      - Tutorial                           ║")
	fmt.Println("║  /personal      - Personal Site                      ║")
	fmt.Println("║  /blog          - Blog                               ║")
	fmt.Println("║  /business      - Business Site                      ║")

	fmt.Println("║                                                      ║")
	fmt.Println("║  Press Ctrl + C to exit                              ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("خطا در اجرای سرور: %v", err)
	}
}

// تابع برای دریافت IP محلی
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "localhost"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// توابع اضافی برای مدیریت فایل‌ها (اختیاری)
func getAllEmbeddedFiles() ([]string, error) {
	var files []string
	err := fs.WalkDir(embeddedFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
