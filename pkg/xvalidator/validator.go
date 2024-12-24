package xvalidator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type CustomValidation struct {
	Tag                      string
	Msg                      string
	Lang                     string
	Func                     validator.Func
	CallValidationEvenIfNull []bool
}

var _customValidations []*CustomValidation

func AddCustomValidations(c ...*CustomValidation) {
	_customValidations = append(_customValidations, c...)
}

func RegisterValidations() {
	for _, v := range _customValidations {
		if v.Lang == "" {
			v.Lang = _defaultLang
		}
		_ = _defaultValidateTranslator.validate.RegisterValidation(v.Tag, v.Func, v.CallValidationEvenIfNull...)
		_ = _defaultValidateTranslator.validate.RegisterTranslation(v.Tag, _defaultValidateTranslator.translators[v.Lang], func(ut ut.Translator) error {
			return ut.Add(v.Tag, v.Msg, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(v.Tag, fe.Field())
			return t
		})
	}
}
