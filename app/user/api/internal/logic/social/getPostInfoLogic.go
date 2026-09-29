// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package social

import (
	"context"

	"github/ikun2021/bestpractice/app/user/api/internal/svc"
	"github/ikun2021/bestpractice/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostInfoLogic {
	return &GetPostInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPostInfoLogic) GetPostInfo(req *types.GetPostInfoReq) (resp *types.GetPostInfoResp, err error) {
	// todo: add your logic here and delete this line

	return
}
