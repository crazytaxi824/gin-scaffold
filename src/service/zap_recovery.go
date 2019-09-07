package service

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"src/global"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ZapRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatus(500)

				path := c.Request.URL.Path
				raw := c.Request.URL.RawQuery
				end := time.Now()

				param := logFormatterParams{
					ClientIP:   c.ClientIP(),
					Method:     c.Request.Method,
					StatusCode: c.Writer.Status(),
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

				msg := fmt.Sprintf("%-7s %s | %3d |%13v |%15s ",
					param.Method,
					param.ReqPath,
					param.StatusCode,
					param.Latency,
					param.ClientIP)

				// 判断 broken pipe
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					// add broken pipe error to gin.context.Error
					_ = c.Error(err.(error))

					// error msg
					param.ErrorMsg = c.Errors

					global.Logger.Error(msg, zap.Any("details", param))
					return

				} else {
					// Panic trace
					var panicPair [2]string
					for i := 3; ; i++ {
						pc, file, line, ok := runtime.Caller(i)
						if !ok {
							break
						}

						fn := runtime.FuncForPC(pc)
						panicPair[0] = fn.Name()
						panicPair[1] = file + ":" + strconv.Itoa(line)

						param.PanicTrace = append(param.PanicTrace, panicPair)
					}

					// Panic Msg
					param.PanicMsg = err.(error).Error()

					// error msg
					param.ErrorMsg = c.Errors

					// this defer is for catching global.Logger.Panic below
					defer func() { recover() }()
					global.Logger.Panic(msg, zap.Any("details", param))
				}
			}
		}()
		c.Next()
	}
}
