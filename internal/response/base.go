package response

import (
	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Response(ctx *gin.Context, result Result, data any) {
	ctx.JSON(200, BaseResponse{
		Code: result.Code(),
		Msg:  result.Msg(),
		Data: data,
	})
}

func ResponseSuccess(ctx *gin.Context, data any) {
	Response(ctx, Success, data)
}

func ResponseFail(ctx *gin.Context, data any) {
	Response(ctx, Fail, data)
}
