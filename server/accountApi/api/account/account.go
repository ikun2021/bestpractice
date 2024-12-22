package account

import (
	"github.com/gin-gonic/gin"
	"github/lunxun9527/bestpractice/pkg/xgin"
	accountApiModel "github/lunxun9527/bestpractice/server/accountApi/model/account"
	accountService "github/lunxun9527/bestpractice/server/accountApi/service/account"
)

var AccountApi = &accountApi{}

type accountApi struct{}

func (*accountApi) GetAccountInfo(c *gin.Context) {
	accountId := c.GetString("accountId")
	resp, err := accountService.UserService.GetAccountInfo(c, &accountApiModel.GetUserInfoReq{AccountId: accountId})
	xgin.ResponseWithLang(c, resp, err)
}
func (*accountApi) AddAccount(c *gin.Context) {
	accountId := c.GetString("accountId")
	resp, err := accountService.UserService.GetAccountInfo(c, &accountApiModel.GetUserInfoReq{AccountId: accountId})
	xgin.ResponseWithLang(c, resp, err)
}
