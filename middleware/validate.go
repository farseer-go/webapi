package middleware

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi/context"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

var validate *validator.Validate
var trans ut.Translator

type Validate struct {
	context.IMiddleware
}

// InitValidate 初始化字段验证
func InitValidate() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	trans, _ = ut.New(zh.New()).GetTranslator("zh")
	_ = zhTranslations.RegisterDefaultTranslations(validate, trans)
	//注册一个函数，获取struct tag里自定义的label作为字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		return name
	})
}

func (receiver *Validate) Invoke(httpContext *context.HttpContext) {
	// 验证dto
	if httpContext.Route.RequestParamIsModel {
		lstError := collections.NewList[string]()
		err := validate.Struct(httpContext.Request.Params[0].Interface())
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			for _, validationError := range validationErrors {
				lstError.Add(validationError.Translate(trans))
			}
		}
		if lstError.Count() > 0 {
			exception.ThrowWebException(403, lstError.ToString(","))
			return
		}
	}
	receiver.IMiddleware.Invoke(httpContext)
}
