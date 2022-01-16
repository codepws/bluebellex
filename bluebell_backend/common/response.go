package common

import (
	"bluebell_backend/common/errno"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseError(ctx *gin.Context, c errno.Error) {

	ctx.JSON(http.StatusOK, c)
}

func ResponseErrorWithMsg(ctx *gin.Context, code int, msg string, data interface{}) {

	errno := errno.NewError(code, msg)
	errno.WithData(data)

	ctx.JSON(http.StatusOK, errno)
}

// 成功响应，不带msg形参，默认success
func ResponseSuccess(ctx *gin.Context, data interface{}) {

	ok := errno.Success.WithData(data)
	ctx.JSON(http.StatusOK, ok)

}

// 成功响应，带msg形参，msg为空设为success
func ResponseSuccessMsg(ctx *gin.Context, msg string, data interface{}) {
	success := errno.NewError(0, msg).WithData(data)
	ctx.JSON(http.StatusOK, success)
}

// 模板响应
func ResponseHtml(c *gin.Context, path string, data interface{}) {
	c.HTML(http.StatusOK, path, data)
}

func Response(ctx *gin.Context, httpStatus int, code int, msg string, data interface{}) {
	errno := errno.NewError(code, msg)
	errno.WithData(data)

	ctx.JSON(httpStatus, errno)
}
