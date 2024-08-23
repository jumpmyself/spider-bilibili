package router

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"spider-bilibili/app/middleware"
	"spider-bilibili/app/schedule/anime"
	"spider-bilibili/app/schedule/documentary"
	"spider-bilibili/app/schedule/guochuang"
	"spider-bilibili/app/schedule/movie"
	"spider-bilibili/app/schedule/tv"
	"spider-bilibili/app/schedule/variety"
)

func Router() {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Use(middleware.LogMiddleware())

	r.GET("/fanju", anime.GetInfo)
	r.GET("/movie", movie.GetInfo)
	r.GET("/guochuang", guochuang.GetInfo)
	r.GET("/tv", tv.GetInfo)
	r.GET("/zhongyi", variety.GetInfo)
	r.GET("/documentary", documentary.GetInfo)

	r.GET("/getimage", func(c *gin.Context) {
		GetImage(c)
	})

	if err := r.Run(":8088"); err != nil {
		panic("gin启动失败")
	}
}

// GetImage 从指定文件夹获取图片
func GetImage(c *gin.Context) {
	// 设置图片文件夹的路径
	imageDir := "./qrcodes" // 只允许从这个文件夹读取图片

	// 从请求中获取图片名称
	imageName := c.Query("imageName")

	// 构建完整的文件路径
	imagePath := filepath.Join(imageDir, imageName)

	// 检查文件是否存在
	if _, err := ioutil.ReadFile(imagePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片未找到"})
		return
	}

	// 读取图片文件
	file, err := ioutil.ReadFile(imagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取图片失败"})
		return
	}

	// 写入响应体
	c.Data(http.StatusOK, "image/png", file)
}
