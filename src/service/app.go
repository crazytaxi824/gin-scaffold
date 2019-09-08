package service

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"src/global"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
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

	svr = http.Server{
		Addr:              global.Config.Service.IP + ":" + strconv.Itoa(global.Config.Service.Port),
		Handler:           App,              // 调度器
		ReadTimeout:       10 * time.Second, // 读取超时
		WriteTimeout:      10 * time.Second, // 响应超时
		IdleTimeout:       10 * time.Second, // 连接空闲超时
		ReadHeaderTimeout: 10 * time.Second, // http header读取超时
	}
}

// 框架中间件
// 注意，Recovery() 放在最外层（最上层）
// 因为middle ware 的执行顺序和 c.Next() 的问题
func frameworkMiddleWare() {
	// gin recovery
	App.Use(ZapRecovery())

	// zap logger - custom logger
	App.Use(ZapLogger())

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
		ctx.AbortWithStatus(405)
	})

	// 404
	App.NoRoute(func(ctx *gin.Context) {
		ctx.AbortWithStatus(404)
	})

	// 500 use abortWithError / ctx.AbortWithError
	// it will use the logger to print the error msg
}

// gin recovery
func setAppRecovery() error {
	errorFile, err := os.OpenFile("frame_err.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	App.Use(gin.RecoveryWithWriter(io.MultiWriter(errorFile, os.Stderr)))
	return nil
}

// gin logger
func setAppLogger() error {
	logFile, err := os.OpenFile("err.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	App.Use(gin.LoggerWithWriter(io.MultiWriter(logFile, os.Stdout)))
	return nil
}
