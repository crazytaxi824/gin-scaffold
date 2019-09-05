package service

import (
	"fmt"
	"src/global"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// logger record
type logFormatterParams struct {
	StatusCode int           `json:"statusCode,omitempty"` // 返回 http code
	TimeStamp  int64         `json:"timeStamp,omitempty"`  // log 时间
	Latency    time.Duration `json:"latency,omitempty"`    // 请求处理时长
	ClientIP   string        `json:"clientIp,omitempty"`   // 请求IP
	Method     string        `json:"method,omitempty"`     // 请求方法
	ReqPath    string        `json:"reqPath,omitempty"`    // 请求路径
	BodySize   int           `json:"bodySize,omitempty"`   // 返回的数据大小

	ErrorMsg []string    `json:"errorMsg,omitempty"` // error 错误信息
	ErrFile  interface{} `json:"errFile,omitempty"`  // error 报错文件

	PanicMsg   string      `json:"panicMsg,omitempty"`   // panic 错误信息
	PanicTrace [][2]string `json:"panicTrace,omitempty"` // panic 错误跟踪
}

func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		end := time.Now()

		// Log only when path is not being skipped
		param := logFormatterParams{
			ClientIP:   c.ClientIP(),
			Method:     c.Request.Method,
			StatusCode: c.Writer.Status(),
			ErrorMsg:   c.Errors.ByType(gin.ErrorTypePrivate).Errors(),
			BodySize:   c.Writer.Size(),
			TimeStamp:  end.Unix(),
		}

		// latency
		param.Latency = end.Sub(start)

		// request path
		if raw != "" {
			path = path + "?" + raw
		}

		param.ReqPath = path

		// error file
		param.ErrFile, _ = c.Get("file")

		msg := fmt.Sprintf("%-7s %s | %3d |%13v |%15s ",
			param.Method,
			param.ReqPath,
			param.StatusCode,
			param.Latency,
			param.ClientIP)

		if param.StatusCode >= 500 {
			global.Logger.Error(msg, zap.Any("details", param))
		} else {
			global.Logger.Info(msg, zap.Any("details", param))
		}
	}
}
