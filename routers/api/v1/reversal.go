package v1

import (
	"context"
	"encoding/hex"
	"net/http"
	"planetd/pkg/e"
	"planetd/pkg/util/format"

	"planetd/pkg/client/planet"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
)

func Reversal(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "request_id", requestID)
	if _rawData, err := c.GetRawData(); err != nil {
		logger.Debugf(ctx, "req=%s", string(_rawData))
	}
	var requestData format.RequestData
	ext := make(map[string]interface{})
	if err := c.ShouldBindJSON(&requestData); err != nil {
		logger.Errorf(ctx, "bind json error: %s", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"respcd":  e.PARAMERR,
			"respmsg": e.GetMsg(e.PARAMERR),
			"resperr": e.GetMsg(e.PARAMERR),
			"chnlsn":  "",
			"pos_ext": ext,
		})
		return
	}
	logger.Debugf(ctx, "requestData: %+v", requestData)
	ups, err := planet.DoReversal(ctx, &requestData)
	if err != nil {
		logger.Errorf(ctx, "do transaction error: %s", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"respcd":  e.DATAERR,
			"respmsg": e.GetMsg(e.DATAERR),
			"resperr": e.GetMsg(e.DATAERR),
			"chnlsn":  "",
			"pos_ext": ext,
		})
		return
	}
	bcdICCData := hex.EncodeToString(ups.ICCSystemRelatedData)
	ext["icccondcode"] = ups.AuthorizationIDResponse
	ext["iccdata"] = bcdICCData
	ext["channel_resp_code"] = ups.ResponseCode

	if ups.ResponseCode != "00" {
		code := "10" + ups.ResponseCode
		msg := e.GetMsg(code)
		c.JSON(http.StatusOK, gin.H{
			"respcd":  code,
			"respmsg": msg,
			"resperr": msg,
			"chnlsn":  ups.RetrievalReferenceNumber,
			"pos_ext": ext,
		})
	} else {
		msg := e.GetMsg(e.SUCCESS)
		c.JSON(http.StatusOK, gin.H{
			"respcd":  e.SUCCESS,
			"respmsg": msg,
			"resperr": msg,
			"chnlsn":  ups.RetrievalReferenceNumber,
			"pos_ext": ext,
		})
		return
	}

}
