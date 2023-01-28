package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"urls/configs"
	_ "urls/docs"
	"urls/handler"
)

func init() {
	configs.InitViper()
}

// @title           Swagger Example API
// @version         1.0
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      mgmo.tech
// @BasePath  /api/v1
func main() {
	configs.InitViper()
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		url := v1.Group("/url")
		{
			url.POST("/shorten", handler.ShortenURL)
			url.GET("/expand/:shorten", handler.ExpandURL)
			url.GET("/hits/:shorten", handler.GetHits)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	_ = r.Run(":8090")
}
