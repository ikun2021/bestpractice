package router

import (
	"github.com/gin-gonic/gin"
	"github/lunxun9527/bestpractice/server/accountApi/api/account"
	"github/lunxun9527/bestpractice/server/accountApi/middleware"
)

func initAccountRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("account")
	userRouter.Use(middleware.TokenValidator())
	{
		userRouter.GET("getAccountInfo", account.AccountApi.GetAccountInfo) // 获取用户
	}
}
