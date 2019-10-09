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

type logFormatterPanicParams struct {
	logFormatterParams

	Panic struct {
		PanicMsg   string       `json:"panicMsg,omitempty"`   // panic 错误信息, 不能用interface 否则无法被 zap.Any 打印
		PanicTrace []panicTrace `json:"panicTrace,omitempty"` // panic 错误跟踪
	} `json:"panic,omitempty"`
}

type panicTrace struct {
	Function string `json:"function,omitempty"` // 错误函数
	Trace    string `json:"trace,omitempty"`    // 错误文件追踪
}

func ZapRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatus(500)

				// 判断 broken pipe
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				brokenPipe := checkBrokenPipe(c, err)

				// format logger msg
				msg, param := formatParam(c, start)

				if brokenPipe {
					global.Logger.Error(msg, zap.Any("details", param))
					return
				} else {
					panicParam := logFormatterPanicParams{
						logFormatterParams: param,
						//Panic: panicDetails,
					}

					// Panic Msg
					// 这里可能为 string，可能为error，所以使用fmt包来自动处理
					panicParam.Panic.PanicMsg = fmt.Sprintf("%s", err)

					// Panic trace
					var panicPair panicTrace
					for i := 3; ; i++ {
						pc, file, line, ok := runtime.Caller(i)
						if !ok {
							break
						}

						// get func name
						fn := runtime.FuncForPC(pc)
						panicPair.Function = fn.Name()

						// get file + line
						panicPair.Trace = file + ":" + strconv.Itoa(line)

						// append to panic trace
						panicParam.Panic.PanicTrace = append(panicParam.Panic.PanicTrace, panicPair)
					}

					// this defer is for catching global.Logger.Panic below
					// 因为 logger.panic 会触发 Panic
					defer func() { recover() }()
					global.Logger.Panic(msg, zap.Any("details", panicParam))
				}
			}
		}()
		// 传递 context
		c.Next()
	}
}

// 检查 是否 broken pipe error。
// 如果是，将 broken pipe 的错误放到 ctx.error 里面，不用 panic 打印。
// true 代表是 broken pipe error
func checkBrokenPipe(ctx *gin.Context, err interface{}) bool {
	var brokenPipe bool
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				brokenPipe = true

				// add broken pipe error to gin.context.Error
				_ = ctx.Error(err.(error))
			}
		}
	}
	return brokenPipe
}
