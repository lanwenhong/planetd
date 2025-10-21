package v1

import (
	"context"
	"net/http"
	"planetd/pkg/e"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
)

func Refund(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)
	logger.Debugf(ctx, "refund")
	msg := e.GetMsg(e.SUCCESS)
	c.JSON(http.StatusOK, gin.H{
		"respcd":  e.SUCCESS,
		"respmsg": msg,
		"resperr": msg,
	})
}
