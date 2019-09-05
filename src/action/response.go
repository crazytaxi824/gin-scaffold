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
func abortWithError(ctx *gin.Context, errMsg error) {
	ctx.AbortWithError(500, errMsg)
	_, file, line, _ := runtime.Caller(1)
	caller := file + ":" + strconv.Itoa(line)
	var files []string
	fileInterface, ok := ctx.Get("file")
	if ok {
		files = fileInterface.([]string)
		files = append(files, caller)
		ctx.Set("file", files)
	} else {
		ctx.Set("file", []string{caller})
	}
}
