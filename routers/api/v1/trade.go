package v1

import (
	"net/http"
	"planetd/pkg/e"

	"github.com/gin-gonic/gin"
)

func Trade(c *gin.Context) {
	data := make(map[string]interface{})
	c.JSON(http.StatusOK, gin.H{
		"respcd":  e.SUCCESS,
		"respmsg": e.GetMsg(e.SUCCESS),
		"resperr": e.GetMsg(e.SUCCESS),
		"data":    data,
	})
}
