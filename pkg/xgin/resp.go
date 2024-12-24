package xgin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
	"github/lunxun9527/bestpractice/common/errs"
	"github/lunxun9527/bestpractice/pkg/i18n"
	"github/lunxun9527/bestpractice/pkg/xvalidator"
	"google.golang.org/grpc/status"
	"net/http"
)

type CommonResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

var (
	Empty = struct{}{}
)

func FailWithLangError(c *gin.Context, err error) {
	ResponseWithLang(c, struct{}{}, err)
}
func FailWithLang(c *gin.Context) {
	FailWithLangError(c, fmt.Errorf("unknown error"))
}
func ResponseWithLang(c *gin.Context, resp interface{}, err error) {
	lang := c.GetHeader("lang")
	if err != nil {
		// 参数校验错误
		var e1 validator.ValidationErrors
		ok := errors.As(err, &e1)
		paramMsg := ""
		if ok {
			paramMsg = ": " + xvalidator.TranslateFirst(lang, err)
			err = errs.ParamValidateFailedErr
		}
		//grpc 以及业务错误
		code := status.Code(err)
		msg := i18n.Translate(lang, cast.ToString(uint32(code)))
		if ok {
			msg += paramMsg
		}
		Result(cast.ToInt(uint32(code)), Empty, msg, c)
	} else {
		Result(SUCCESS, resp, "success", c)
	}
}

func Response(c *gin.Context, resp interface{}, err error) {
	if err != nil {
		r, _ := status.FromError(err)
		Result(cast.ToInt(r.Code()), Empty, r.Message(), c)
	} else {
		Result(SUCCESS, resp, "success", c)
	}
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, CommonResponse{
		code,
		msg,
		data,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "success", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}
