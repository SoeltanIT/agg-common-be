package common

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/pkg/errors"
)

var (
	translator *ut.UniversalTranslator
	validate   *validator.Validate
)

func init() {
	var err error
	validate, err = registerTranslator()
	if nil != err {
		panic(err)
	}
}

// Validator returns the custom validator instance
func Validator() *validator.Validate {
	return validate
}

// registerTranslator registers the custom translator for validator
func registerTranslator() (*validator.Validate, error) {
	v := validator.New()

	english := en.New()
	translator = ut.New(english, english)

	translatorEn, ok := translator.GetTranslator("en")
	if !ok {
		return nil, errors.New("cannot found message translator")
	}

	if err := enTranslations.RegisterDefaultTranslations(v, translatorEn); nil != err {
		return nil, errors.Wrap(err, "registering default translation")
	}

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v, nil
}
