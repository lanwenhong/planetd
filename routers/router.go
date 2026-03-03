package routers

import (
	"net/http"
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
		apiV1.POST("/swipe", v1.Trade)
		apiV1.POST("/tip", v1.TradeTip)
		apiV1.POST("/refund", v1.Trade)
		apiV1.POST("/close", v1.Trade)
		apiV1.POST("/void_refund", v1.Trade)
		apiV1.POST("/reversal", v1.Reversal)
		apiV1.POST("/key_exchange", v1.KeyExchange)
	}
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	return r
}
