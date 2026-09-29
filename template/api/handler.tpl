package {{.PkgName}}

import (
	"net/http"
    "gitee.com/lilei12138/xpkg/response"
    {{if .HasRequest}}
    "gitee.com/lilei12138/xpkg/xerrs"
	"github.com/zeromicro/go-zero/rest/httpx"
	 {{end}}
	{{.ImportPackages}}
)

{{if .HasDoc}}{{.Doc}}{{end}}
func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
			response.GoZeroRespWithLang(w, r, nil, xerrs.ParamValidateFailedErr.WarpMessage(err.Error()))
			return
		}

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		response.GoZeroRespWithLang(w, r, resp, err)

	}
}
