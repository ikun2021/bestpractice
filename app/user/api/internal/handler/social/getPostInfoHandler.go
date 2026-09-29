package social

import (
	"gitee.com/lilei12138/xpkg/response"
	"net/http"

	"gitee.com/lilei12138/xpkg/xerrs"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github/ikun2021/bestpractice/app/user/api/internal/logic/social"
	"github/ikun2021/bestpractice/app/user/api/internal/svc"
	"github/ikun2021/bestpractice/app/user/api/internal/types"
)

func GetPostInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPostInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			response.GoZeroRespWithLang(w, r, nil, xerrs.ParamValidateFailedErr.WarpMessage(err.Error()))
			return
		}

		l := social.NewGetPostInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetPostInfo(&req)
		response.GoZeroRespWithLang(w, r, resp, err)

	}
}
