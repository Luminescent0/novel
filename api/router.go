package api

import (
	"github.com/gin-gonic/gin"
)

func InitEngine() {
	engine := gin.Default()
	//engine.Use(Cors()) //跨域

	//oauth

	engine.POST("/login", login)       //登录
	engine.POST("/register", register) //注册

	userGroup := engine.Group("/user")
	{
		userGroup.Use(JwtAuthMiddleware)
		userGroup.POST("/password", changePassword)
		userGroup.POST("/liked", liked)

	}
	bookGroup := engine.Group("/book")
	{
		bookGroup.GET("/search", search)
	}
	execGroup := engine.Group("/exec")
	{
		execGroup.Use(JwtAuthMiddleware)

	}
	err := engine.Run(":8080")
	if err != nil {
		return
	}
}
