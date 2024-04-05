package tool

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespErrorWithData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"info": data,
	})
}

func RespInternalError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"info": "服务器错误",
	})
}

func RespSuccessful(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"info": "成功",
	})
}

func RespSuccessfulWithData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"info": "成功",
		"data": data,
	})
}
func RespSuccessfulLogin(ctx *gin.Context, refresh_token, token string) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success",
		"data": gin.H{
			"refresh_token": refresh_token,
			"token":         token,
		},
	})
}
func RespAbortWithStatus(ctx *gin.Context, data string) {
	ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"code": -1,
		"msg":  data,
	})
}
