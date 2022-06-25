package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type customValidator struct {
	validator *validator.Validate
}

// map校验类型
type MapValidate struct {
	Value map[string]interface{}
	Rules map[string]interface{}
}

func NewCustomValidator() *customValidator {
	return &customValidator{validator: validator.New()}
}

// Validate implements echo.Validator 暂时只校验map格式数据
func (cv *customValidator) Validate(i interface{}) error {
	data := i.(*MapValidate)
	errMap := cv.validator.ValidateMap(
		data.Value,
		data.Rules,
	)

	fmt.Printf("validate:%#v, rules: %#v, errMap:%#v \n", data.Value, data.Rules, errMap)

	fmt.Println("errMapLen,", len(errMap))

	if len(errMap) == 0 {
		return nil
	} else {
		msgError := ""
		for k, v := range errMap {
			if msgError == "" {
				msgError = fmt.Sprintf("%s(%s)", k, v)
			} else {
				msgError = fmt.Sprintf("%s,%s(%s)", msgError, k, v)
			}
		}
		return fmt.Errorf("validate error msg: %s", msgError)
	}
}
