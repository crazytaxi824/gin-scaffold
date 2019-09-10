package action

import (
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 响应单键值的数据
type RespData struct {
	Error string      `json:"error"`
	Data  interface{} `json:"data,omitempty"`
}

// 响应多键值的数据
type RespMapData struct {
	Error string                 `json:"error"`
	Data  map[string]interface{} `json:"data,omitempty"`
}

// 构建一个RespMapData
func makeRespMapData() (respData RespMapData) {
	respData.Data = make(map[string]interface{})
	return
}

// 返回 500 错误，追踪 caller
func abortWithError(ctx *gin.Context, err error) {
	// 追加错误行数到 gin.Error 中，使用 meta 属性
	_, file, line, _ := runtime.Caller(1)
	ctx.AbortWithError(500, err).SetMeta(file + ":" + strconv.Itoa(line))
}
