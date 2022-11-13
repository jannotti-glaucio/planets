package validations

import (
	"time"

	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/core/tools/validations/validator"

	"github.com/asaskevich/govalidator"
)

var simpleValidator = map[string]struct {
	validator func(str string) bool
	message   string
	code      int
}{
	"isRequired": {
		validator: govalidator.IsNotNull,
		message:   communication.New().Mapping["validate_required"].Message,
		code:      communication.New().Mapping["validate_required"].Code,
	},
	"isRequiredCreated": {
		validator: govalidator.IsNotNull,
		message:   communication.New().Mapping["validate_required"].Message,
		code:      communication.New().Mapping["validate_required"].Code,
	},
	"isRequiredUpdated": {
		validator: govalidator.IsNotNull,
		message:   communication.New().Mapping["validate_required"].Message,
		code:      communication.New().Mapping["validate_required"].Code,
	},
	"isEmail": {
		validator: govalidator.IsEmail,
		message:   communication.New().Mapping["validate_invalid"].Message,
		code:      communication.New().Mapping["validate_invalid"].Code,
	},
}

var datesValidator = map[string]struct {
	validator  func(str string) (bool, time.Time)
	identifier string
	message    string
	code       int
}{
	"isDate": {
		validator: validator.IsDate,
		message:   communication.New().Mapping["validate_date"].Message,
		code:      communication.New().Mapping["validate_date"].Code,
	},
}
