package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"novel/service"
	"novel/tool"
)

func search(ctx *gin.Context) {
	bookName := ctx.PostForm("bookName")
	novel, err := service.FindNovel(bookName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			tool.RespErrorWithData(ctx, "该小说不存在")
			return
		}
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessfulWithData(ctx, novel)
}
