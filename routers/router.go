package routers

import (
	"planetd/middleware"
	"planetd/pkg/setting"
	v1 "planetd/routers/api/v1"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.LoogerToFile())
	gin.SetMode(setting.Conf.RunMode)
	apiV1 := r.Group("/trade/v1")
	{
		apiV1.POST("/payment", v1.Trade)
	}
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r
}
