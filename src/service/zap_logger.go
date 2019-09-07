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

	ErrorMsg []*gin.Error `json:"errors,omitempty"` // error 错误信息,包含error 和 path/meta
}

func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// format logger msg
		msg, param := formatParam(c, start)

		if param.StatusCode >= 500 {
			global.Logger.Error(msg, zap.Any("details", param))
		} else {
			global.Logger.Info(msg, zap.Any("details", param))
		}
	}
}

func formatParam(c *gin.Context, start time.Time) (msg string, param logFormatterParams) {
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	end := time.Now()
	param = logFormatterParams{
		ClientIP:   c.ClientIP(),
		Method:     c.Request.Method,
		StatusCode: c.Writer.Status(),
		BodySize:   c.Writer.Size(),
		TimeStamp:  end.Unix(),
		Latency:    end.Sub(start),
		ErrorMsg:   c.Errors,
	}

	// request path
	if raw != "" {
		path = path + "?" + raw
	}

	param.ReqPath = path

	msg = fmt.Sprintf("%-7s %s | %3d |%13v |%15s ",
		param.Method,
		param.ReqPath,
		param.StatusCode,
		param.Latency,
		param.ClientIP)

	return
}
