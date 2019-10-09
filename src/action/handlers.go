package action

import (
	"errors"
	"log"
	"src/global"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetHandler(ctx *gin.Context) {
	// 获取
	log.Println(ctx.Query("id"))

	log.Println(ctx.Params.ByName("id"))
	log.Println(ctx.Params.Get("id"))
	log.Println(ctx.Param("id"))

	data := make(map[string]interface{})
	data["a"] = 123

	log.Println("------------------------------------")
	ctx.JSON(200, data)
}

func PostHandler(ctx *gin.Context) {
	// 获取midware传参
	log.Println(ctx.Get("mid"))
	// 获取参数
	log.Println(ctx.Query("id"))
	log.Println(ctx.PostForm("name"))
	log.Println(ctx.PostForm("age"))

	//panic("hello world")
	_, err := strconv.ParseInt("1110", 10, 64)
	if err != nil {
		ctx.AbortWithError(500, err)
		global.Logger.Sugar().Error(err)
		return
	}

	number := 123
	str := "hello"

	ctx.String(200, "post | %s | %d", str, number)
	return
}

func PanicTest(ctx *gin.Context) {
	// 先添加两个错误，再触发panic
	_, err := strconv.Atoi("haha")
	if err != nil {
		abortWithError(ctx, err)
	}

	_, err = strconv.Atoi("xixi")
	if err != nil {
		abortWithError(ctx, err)
	}

	var is []int
	// 触发 panic
	log.Println(is[2])
}

func ErrTest(ctx *gin.Context) {
	abortWithError(ctx, errors.New("new error1"))
	abortWithError(ctx, errors.New("new error2"))
	return
}
