package account

import (
	"github.com/go-playground/validator/v10"
	"github/lunxun9527/bestpractice/pkg/xvalidator"
)

type GetAccountInfoReq struct {
	AccountId string `form:"accountId" binding:"required"`
}

type RegisterUserReq struct{}

type AddAccountReq struct {
	AccountID   string `json:"accountId" binding:"required"`
	AccountName string `json:"accountName" binding:"required,min=6"`
	Password    string `json:"password" binding:"account_password"`
}

func init() {
	xvalidator.AddCustomValidations(&xvalidator.CustomValidation{
		Tag:  "account_password",
		Lang: "zh",
		Msg:  "{0} 长度需在8到64个字符串之间",
		Func: func(fl validator.FieldLevel) bool {
			// 检查长度是否在 8 到 64 字节之间
			if len(fl.Field().String()) < 8 || len(fl.Field().String()) > 64 {
				return false
			}
			return true
		},
		CallValidationEvenIfNull: nil,
	}, &xvalidator.CustomValidation{
		Tag:  "account_password",
		Lang: "en",
		Msg:  "{0} the length is between 8 and 64 bytes",
		Func: func(fl validator.FieldLevel) bool {
			// 检查长度是否在 8 到 64 字节之间
			if len(fl.Field().String()) < 8 || len(fl.Field().String()) > 64 {
				return false
			}
			return true
		},
		CallValidationEvenIfNull: nil,
	})

}
