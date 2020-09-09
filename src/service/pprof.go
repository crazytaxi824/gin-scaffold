package service

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

// 在 gin 的路由中注册 pprof，gin logger 会打印访问记录
func pprofRegister(engine *gin.Engine, prefixPath ...string) {
	var p string

	if prefixPath == nil {
		p = "/debug/pprof"
	} else {
		p = prefixPath[0]
	}

	{
		// pprof
		prefixRouter := engine.Group(p)
		prefixRouter.GET("/", pprofHandler(pprof.Index))
		prefixRouter.GET("/cmdline", pprofHandler(pprof.Cmdline))
		prefixRouter.GET("/profile", pprofHandler(pprof.Profile))
		prefixRouter.POST("/symbol", pprofHandler(pprof.Symbol))
		prefixRouter.GET("/symbol", pprofHandler(pprof.Symbol))
		prefixRouter.GET("/trace", pprofHandler(pprof.Trace))
		prefixRouter.GET("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
		prefixRouter.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
		prefixRouter.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
		prefixRouter.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
		prefixRouter.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
		prefixRouter.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
	}
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// 使用新的 mux 添加 pprof，在 gin 路由之外，gin logger 不会打印访问记录。
func pprofMux() http.Handler {
	pprofMux := http.NewServeMux()
	pprofMux.HandleFunc(pprofPath, pprof.Index)
	pprofMux.HandleFunc(pprofPath+"cmdline", pprof.Cmdline)
	pprofMux.HandleFunc(pprofPath+"profile", pprof.Profile)
	pprofMux.HandleFunc(pprofPath+"symbol", pprof.Symbol)
	pprofMux.HandleFunc(pprofPath+"trace", pprof.Trace)
	pprofMux.HandleFunc(pprofPath+"allocs", pprof.Handler("allocs").ServeHTTP)
	pprofMux.HandleFunc(pprofPath+"block", pprof.Handler("block").ServeHTTP)
	pprofMux.HandleFunc(pprofPath+"goroutine", pprof.Handler("goroutine").ServeHTTP)
	pprofMux.HandleFunc(pprofPath+"heap", pprof.Handler("heap").ServeHTTP)
	pprofMux.HandleFunc(pprofPath+"mutex", pprof.Handler("mutex").ServeHTTP)
	pprofMux.HandleFunc(pprofPath+"threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	return pprofMux
}
