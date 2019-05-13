package Models

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strategy-control-platform/strategy-control-platform-server/Model"
	"tipu.com/go-framework/Code"
)

// Result 返回结果时的通用结构
type Result struct {
	Code int16       `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// ResultList 列表型结果的结构
type ResultList struct {
	Code  int16         `json:"code"`
	Rows  []interface{} `json:"rows"`
	Total int           `json:"total"`
	Msg   string        `json:"msg"`
}

// ResultSuccess 返回结果的标准方式，传递结构体数据
func ResultSuccess(ctx *gin.Context, data interface{}) {
	ResultSuccessWithMsg(ctx, data, "")
}

// ResultSuccess 返回结果的标准方式，传递结构体数据和消息
func ResultSuccessWithMsg(ctx *gin.Context, data interface{}, msg string) {
	ctx.JSON(http.StatusOK, &Model.Result{Code: Code.NORMAL, Data: data, Msg: msg})
}

// ResultListSuccess 返回列表结果，需要传泛型数组和消息
func ResultListSuccess(ctx *gin.Context, data []interface{}, total int) {
	ResultListSuccessWithMsg(ctx, data, total, "")
}

// ResultListSuccess 返回列表结果，需要传泛型数组和消息
func ResultListSuccessWithMsg(ctx *gin.Context, data []interface{}, total int, msg string) {
	ctx.JSON(http.StatusOK, &Model.ResultList{Code: Code.NORMAL, Rows: data, Total: total, Msg: msg})
}

// ResultFail 返回错误结果及错误code
func ResultFail(ctx *gin.Context, code uint16, err error) {
	ctx.JSON(http.StatusOK, &Model.Result{Code: code, Msg: err.Error()})
}
