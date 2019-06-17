package Models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tipu.com/go-framework/Code"
)

// Result 返回结果时的通用结构
type Result struct {
	Code uint16      `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// ResultSuccess 返回结果的标准方式，传递结构体数据
func ResultSuccess(ctx *gin.Context, data interface{}) {
	ResultSuccessWithMsg(ctx, data, "")
}

// ResultSuccess 返回结果的标准方式，传递结构体数据和消息
func ResultSuccessWithMsg(ctx *gin.Context, data interface{}, msg string) {
	ctx.JSON(http.StatusOK, Result{Code: Code.NORMAL, Data: data, Msg: msg})
}

// ResultFail 返回错误结果及错误code
func ResultFail(ctx *gin.Context, code uint16, err error) {
	ctx.JSON(http.StatusOK, Result{Code: code, Msg: err.Error()})
}

func (r Result) String() string {
	return fmt.Sprint("Code: %d, 内容: %s, 数据: %s", r.Code, r.Msg, r.Data)
}
