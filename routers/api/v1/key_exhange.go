package v1

import (
	"context"
	"net/http"
	"planetd/pkg/client/planet"
	"planetd/pkg/e"
	"planetd/pkg/util/format"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
)

func KeyExchange(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)
	logger.Infof(ctx, "key exchange")
	var requestData format.RequestData
	ext := make(map[string]interface{})
	if err := c.ShouldBindJSON(&requestData); err != nil {
		logger.Errorf(ctx, "bind json error: %s", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"respcd":  e.PARAMERR,
			"respmsg": e.GetMsg(e.PARAMERR),
			"resperr": e.GetMsg(e.PARAMERR),
			"chnlsn":  "",
			"ext":     ext,
		})
		return
	}
	logger.Debugf(ctx, "requestData: %+v", requestData)
	ups, err := planet.DoExchange(ctx, &requestData)
	if err != nil {
		logger.Errorf(ctx, "key exchange error: %s", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"respcd":  e.DATAERR,
			"respmsg": e.GetMsg(e.DATAERR),
			"resperr": e.GetMsg(e.DATAERR),
			"chnlsn":  "",
			"ext":     ext,
		})
		return
	}
	msg := e.GetMsg(e.SUCCESS)
	c.JSON(http.StatusOK, gin.H{
		"respcd":  e.SUCCESS,
		"respmsg": msg,
		"resperr": msg,
		"chnlsn":  ups.RetrievalReferenceNumber,
		"ext":     ext,
	})
}
