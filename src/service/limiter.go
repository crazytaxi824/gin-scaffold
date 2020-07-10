package service

import (
	"net/http"

	"src/global"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiter *rate.Limiter

func SetLimiter() {
	// 初始化 limiter
	limiter = rate.NewLimiter(rate.Limit(global.Config.Service.Limiter), global.Config.Service.Limiter)
}

func LimiterMiddle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 超过访问量，直接抛弃
		if !limiter.Allow() {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
	}
}
