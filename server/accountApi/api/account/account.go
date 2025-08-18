package account

import (
	"github.com/gin-gonic/gin"
	"github.com/luxun9527/zlog"
	"github/lunxun9527/bestpractice/pkg/xgin"
	accountApiModel "github/lunxun9527/bestpractice/server/accountApi/model/account"
	accountService "github/lunxun9527/bestpractice/server/accountApi/service/account"
)

var AccountApi = &accountApi{}

type accountApi struct{}

// GetAccountInfo
// @Tags      GetAccountInfo
// @Summary   获取账号信息
// @Security  Account
// @Produce   application/json
// @Param     data  query  accountApiModel.GetAccountInfoReq  true  "请求参数"
// @Success   200   {array}  accountApiModel.GetAccountInfoResp "账号信息"
// @Router    /account/getAccountInfo [get]
func (*accountApi) GetAccountInfo(c *gin.Context) {
	var req accountApiModel.GetAccountInfoReq
	zlog.InfofCtx(c.Request.Context(), "get account info")

	if err := c.ShouldBindQuery(&req); err != nil {
		xgin.FailWithLangError(c, err)
		return
	}
	zlog.InfofCtx(c.Request.Context(), "get account info")
	resp, err := accountService.UserService.GetAccountInfo(c, &accountApiModel.GetAccountInfoReq{AccountId: req.AccountId})
	xgin.ResponseWithLang(c, resp, err)
}
func (*accountApi) AddAccount(c *gin.Context) {
	var req accountApiModel.AddAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xgin.FailWithLangError(c, err)
		return
	}
	if err := accountService.UserService.AddAccount(c, &req); err != nil {
		zlog.Errorf("AddAccount err: %v", err)
		xgin.FailWithLangError(c, err)
		return
	}
	xgin.Ok(c)
}
