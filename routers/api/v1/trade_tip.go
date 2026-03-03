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

func TradeTip(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "request_id", requestID)
	var requestData format.RequestData
	dbExt := make(map[string]interface{})
	ext := make(map[string]interface{})
	ext["icccondcode"] = ""
	ext["iccdata"] = ""
	ext["channel_resp_code"] = ""
	dbExt["tag_fa"] = "F"
	dbExt["tag_tc"] = "5"
	dbExt["processing_cd"] = ""
	if err := c.ShouldBindJSON(&requestData); err != nil {
		logger.Errorf(ctx, "bind json error: %s", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"respcd":  e.PARAMERR,
			"respmsg": e.GetMsg(e.PARAMERR),
			"resperr": e.GetMsg(e.PARAMERR),
			"chnlsn":  "",
			"pos_ext": ext,
			"ext":     dbExt,
		})
		return
	}
	logger.Debugf(ctx, "requestData: %+v", requestData)
	ups, err := planet.DoTipTransaction(ctx, &requestData)
	if err != nil {
		logger.Errorf(ctx, "do transaction error: %s", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"respcd":  e.DATAERR,
			"respmsg": e.GetMsg(e.DATAERR),
			"resperr": e.GetMsg(e.DATAERR),
			"chnlsn":  "",
			"pos_ext": ext,
			"ext":     dbExt,
		})
		return
	}
	dbExt["processing_cd"] = ups.ProcessingCd
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
			"ext":     dbExt,
		})
	} else {
		msg := e.GetMsg(e.SUCCESS)
		c.JSON(http.StatusOK, gin.H{
			"respcd":  e.SUCCESS,
			"respmsg": msg,
			"resperr": msg,
			"chnlsn":  ups.RetrievalReferenceNumber,
			"pos_ext": ext,
			"ext":     dbExt,
		})
		return
	}

}
