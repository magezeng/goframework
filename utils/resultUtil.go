package utils

import (
	DO "go-framework/model/do"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResultSuccess 返回结果的标准方式，传递结构体数据和消息
func ResultSuccess(ctx *gin.Context, data interface{}, msg string) {
	ctx.JSON(http.StatusOK, &DO.Result{Code: http.StatusOK, Data: data, Msg: msg})
}

// ResultListSuccess 返回列表结果，需要传泛型数组和消息
func ResultListSuccess(ctx *gin.Context, data []interface{}, total int, msg string) {
	ctx.JSON(http.StatusOK, &DO.ResultList{Code: http.StatusOK, Rows: data, Total: total, Msg: msg})
}

// ResultFail 返回错误结果
func ResultFail(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, &DO.Result{Code: http.StatusBadRequest, Msg: err.Error()})
}
