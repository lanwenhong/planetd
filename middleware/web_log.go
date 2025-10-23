package middleware

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
)

func LoogerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		handleStart := time.Now()
		requestID := util.GenerateUniqueStringWithTimestamp("id")
		c.Request.Header.Set("X-Request-ID", requestID)
		ctx := context.WithValue(context.Background(), "trace_id", requestID)

		var reqData []byte
		if c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPost {
			bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
			logger.Debugf(ctx, "bodyBytes: %s", string(bodyBytes))
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			if len(bodyBytes) >= 1024 {
				reqData = bodyBytes[0:1024]
			} else {
				reqData = bodyBytes
			}
		}

		// 创建自定义的ResponseWriter来捕获响应数据
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		reqMethod := c.Request.Method
		// 状态码
		statusCode := c.Writer.Status()

		path := c.Request.URL.Path

		rawQuery := c.Request.URL.RawQuery

		if rawQuery == "" {
			rawQuery = "-"
		}

		// 获取响应数据
		var respData []byte
		if blw.body.Len() > 0 {
			respBytes := blw.body.Bytes()
			if len(respBytes) >= 1024 {
				respData = respBytes[0:1024]
			} else {
				respData = respBytes
			}
			logger.Debugf(ctx, "response body: %s", string(respData))
		}

		// 请求IP
		clientIp := c.ClientIP()
		handleEnd := time.Now()
		handleTime := handleEnd.Sub(handleStart)
		logger.Infof(ctx, "%d|%v|%v|%s|%s|%s|%s|%s|%s",
			statusCode,
			latencyTime,
			handleTime,
			clientIp,
			reqMethod,
			path,
			rawQuery,
			reqData,
			respData,
		)
	}
}

// bodyLogWriter 自定义ResponseWriter来捕获响应数据
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
