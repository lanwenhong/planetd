package v1

import (
	"context"
	"net/http"
	"planetd/pkg/e"
	"planetd/pkg/util/format"

	"planetd/pkg/client/planet"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
)

func Trade(c *gin.Context) {
	data := make(map[string]interface{})
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)
	var requestData format.RequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		logger.Errorf(ctx, "bind json error: %s", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"respcd":  e.PARAMERR,
			"respmsg": e.GetMsg(e.PARAMERR),
			"resperr": e.GetMsg(e.PARAMERR),
			"data":    data,
		})
		return
	}
	logger.Debugf(ctx, "requestData: %+v", requestData)
	ups, err := planet.DoTransaction(ctx, &requestData)
	if err != nil {
		logger.Errorf(ctx, "do transaction error: %s", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"respcd":  e.DATAERR,
			"respmsg": e.GetMsg(e.DATAERR),
			"resperr": e.GetMsg(e.DATAERR),
			"data":    data,
		})
		return
	}
	if ups.ResponseCode != "00" {
		code := "10" + ups.ResponseCode
		msg := e.GetMsg(code)
		c.JSON(http.StatusOK, gin.H{
			"respcd":  code,
			"respmsg": msg,
			"resperr": msg,
			"data":    data,
		})
	} else {
		msg := e.GetMsg(e.SUCCESS)
		c.JSON(http.StatusOK, gin.H{
			"respcd":  e.SUCCESS,
			"respmsg": msg,
			"resperr": msg,
			"data":    data,
		})
		return
	}

}
