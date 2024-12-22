package account

type GetUserInfoReq struct {
	AccountId string `form:"accountId"`
}

type RegisterUserReq struct{}

type AddAccountReq struct {
	AccountID   string `json:"account_id" binding:"required"`
	AccountName string `json:"account_name" binding:"required,min=6"`
	Password    string `json:"password" binding:"required,min=6"`
}
