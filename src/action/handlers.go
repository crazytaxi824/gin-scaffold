package action

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetHandler(ctx *gin.Context) {
	// 获取
	log.Println(ctx.Query("id"))

	// log.Println(ctx.Params.ByName("id"))
	// log.Println(ctx.Params.Get("id"))
	log.Println(ctx.Param("id"))

	// 返回数据
	data := make(map[string]interface{})
	data["a"] = 123

	ctx.JSON(http.StatusOK, data)
}

func PostHandler(ctx *gin.Context) {
	// 获取 midware 传参
	log.Println(ctx.Get("mid"))
	// 获取参数
	log.Println(ctx.Query("id"))
	log.Println(ctx.PostForm("name"))
	log.Println(ctx.PostForm("age"))

	// panic("hello world")
	_, err := strconv.ParseInt("1110", 10, 64)
	if err != nil {
		abortWithError(ctx, err)
		return
	}

	number := 123
	str := "hello"

	ctx.String(http.StatusOK, "post | %s | %d", str, number)
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
}

// BindingJSONBody Test ShouldBindJSON
//
// GET http://localhost:18080/bind/json
// Content-Type: application/json
// {"usr":"manu", "pwd":"123"}
// ###
func BindingJSONBody(ctx *gin.Context) {
	type login struct {
		User     string `form:"user" json:"usr" xml:"user"  binding:"required"`
		Password string `form:"password" json:"pwd" xml:"password" binding:"required"`
	}

	var usr login
	if err := ctx.ShouldBindJSON(&usr); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(usr)

	ctx.AbortWithStatus(http.StatusNoContent)
}

// BindingQueryGet Test ShouldBindQuery
// GET http://localhost:18080/bind/query?user=abc&password=123
// Content-Type: text/plain
// ###
func BindingQueryGet(ctx *gin.Context) {
	type login struct {
		User     string `form:"user" json:"usr" xml:"user"  binding:"required"`
		Password string `form:"password" json:"pwd" xml:"password" binding:"required"`
	}

	var usr login
	if err := ctx.ShouldBindQuery(&usr); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(usr)

	ctx.AbortWithStatus(http.StatusNoContent)
}

// BindingQueryPost POST 请求使用 ShouldBind 可以自动绑定属性。
// POST http://localhost:18080/bind/post
// Content-Type: application/x-www-form-urlencoded
// user=abc&password=123
// ###
func BindingQueryPost(ctx *gin.Context) {
	type login struct {
		User     string `form:"user" json:"usr" xml:"user"  binding:"required"`
		Password string `form:"password" json:"pwd" xml:"password" binding:"required"`
	}

	var usr login
	if err := ctx.ShouldBind(&usr); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(usr)

	ctx.AbortWithStatus(http.StatusNoContent)
}

// BindingQueryGetAndPost POST请求，同时url有query字段
// 使用 ShouldBind 可以自动绑定属性。
// POST http://localhost:18080/bind/gap?user=abcdefghijk
// Content-Type: application/x-www-form-urlencoded
// password=123
// ###
func BindingQueryGetAndPost(ctx *gin.Context) {
	type login struct {
		User     string `form:"user" json:"usr" xml:"user"  binding:"required,my-val"`
		Password int    `form:"password" json:"pwd" xml:"password" binding:"required"`
		Color    string `form:"color" binding:"required,iscolor"`
	}

	var usr login
	if err := ctx.ShouldBind(&usr); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(usr)

	ctx.AbortWithStatus(http.StatusNoContent)
}
