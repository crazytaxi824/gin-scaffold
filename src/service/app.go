package service

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"src/global"
)

var (
	App *gin.Engine
	svr http.Server
)

func Start() {
	EngineConfig()

	// 在新协程中启动服务，方便实现退出等待
	go func() {
		global.Logger.Info("HTTP服务 " + svr.Addr + " 启动成功")
		if err := svr.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				global.Logger.Info("服务已退出")
			} else {
				global.Logger.Fatal(err.Error())
			}
		}
	}()

	// 退出进程时等待
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit
	// 指定退出超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Config.Service.QuitWaitTimeout)*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		global.Logger.Fatal(err.Error())
	}
}

func EngineConfig() {
	// ReleaseMode
	gin.SetMode(gin.ReleaseMode)

	// gin engine
	App = gin.New()

	// set framework middle ware
	frameworkMiddleWare()

	// deal with 405 404
	eventHandler()

	// set routers
	setRouters()

	const timeout = 10

	svr = http.Server{
		Addr:              global.Config.Service.IP + ":" + strconv.Itoa(global.Config.Service.Port),
		Handler:           App,                   // 调度器
		ReadTimeout:       timeout * time.Second, // 读取超时
		WriteTimeout:      timeout * time.Second, // 响应超时
		IdleTimeout:       timeout * time.Second, // 连接空闲超时
		ReadHeaderTimeout: timeout * time.Second, // http header读取超时
	}
}

// 框架中间件
// 注意，Recovery() 放在最外层（最上层）
// 因为 middle ware 的执行顺序和 c.Next() 的问题
func frameworkMiddleWare() {
	if global.Config.Service.Limiter > 0 {
		App.Use(LimiterMiddle())
	}

	// recovery with zap logger
	App.Use(ZapRecovery())

	// zap logger - custom logger
	App.Use(ZapLogger())

	// allow cross origin
	App.Use(CORS())

	// 是否将 error 信息返回给前端，service debug = true 模式
	// gin.ErrorLogger() 的用途是将 error 的内容返回给前端
	if global.Config.Service.Debug {
		App.Use(gin.ErrorLogger())
	}
}

// 404 405
func eventHandler() {
	// 405
	App.HandleMethodNotAllowed = true
	App.NoMethod(func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusMethodNotAllowed)
	})

	// 404
	App.NoRoute(func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusNotFound)
	})

	// 500 use abortWithError / ctx.AbortWithError
	// it will use the logger to print the error msg
}

// nolint:deadcode,unused // gin default recovery
func setAppRecovery() error {
	// #nosec
	errorFile, err := os.OpenFile("frame_err.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	App.Use(gin.RecoveryWithWriter(io.MultiWriter(errorFile, os.Stderr)))
	return nil
}

// nolint:deadcode,unused // gin default logger
func setAppLogger() error {
	// #nosec
	logFile, err := os.OpenFile("err.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	App.Use(gin.LoggerWithWriter(io.MultiWriter(logFile, os.Stdout)))
	return nil
}
