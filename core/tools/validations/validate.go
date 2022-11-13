package validations

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/core/tools/validations/validator"

	"github.com/asaskevich/govalidator"
)

//ValidateStruct responsible for validate structs
func ValidateStruct(data interface{}, prefix string) (communicationFields []communication.Fields) {
	vl := reflect.ValueOf(data)
	return validateStruct(vl, prefix)
}

func validateStruct(vl reflect.Value, prefix string) (communicationFields []communication.Fields) {

	//must be a pointer
	if vl.Kind() == reflect.Ptr || vl.Kind() == reflect.Struct {

		var elm reflect.Value

		if vl.Kind() == reflect.Ptr {
			elm = vl.Elem()
		} else {
			elm = vl
		}

		for i := 0; i < elm.NumField(); i++ {
			valueField := elm.Field(i)
			typeField := elm.Type().Field(i)
			field := elm.Type().Field(i).Name
			tagValidate := typeField.Tag.Get("validate")
			tagField := field
			prefixField := prefix

			//try to get custom tags
			if tag := elm.Type().Field(i).Tag.Get("json"); tag != "" {
				tagField = tag
			} else if tag := elm.Type().Field(i).Tag.Get("form"); tag != "" {
				tagField = tag
			}

			//Normalise prefix
			if prefixField != "" {
				prefixField = fmt.Sprintf("%s[%s]", prefix, tagField)
			} else {
				prefixField = tagField
			}

			//verify field is public
			if validator.Matches(strings.Split(field, "")[0], "[A-Z]") {
				//If the field type is a struct then validate individual
				if valueField.Kind() == reflect.Struct {
					switch valueField.Interface().(type) {
					case time.Time:
						//Struct time
						validateErrors := validateDate(prefixField, tagValidate, valueField)
						communicationFields = append(communicationFields, validateErrors...)
					default:
						//Simple struct
						validateErrors := validateStruct(valueField, prefixField)
						communicationFields = append(communicationFields, validateErrors...)
					}
				} else
				//If the type is a slice then loop and validate one by one
				if valueField.Kind() == reflect.Slice {
					//Simple slices
					validateErrors := validateSimpleSlice(prefixField, tagValidate, valueField)
					communicationFields = append(communicationFields, validateErrors...)
				} else {
					//Simple Field
					validateErrors := validateField(prefixField, tagValidate, valueField)
					communicationFields = append(communicationFields, validateErrors...)
				}
			}
		}
	}

	return
}

func validateSimpleSlice(prefix string, tagValidate string, slice reflect.Value) (communicationFields []communication.Fields) {
	comm := communication.New()
	if slice.Kind() == reflect.Slice {
		elementsOfSlice := reflect.ValueOf(slice.Interface())
		quantityElements := elementsOfSlice.Len()
		tagsOfSlice := getTags(tagValidate)

		if quantityElements == 0 {
			for x := range tagsOfSlice {
				if tagsOfSlice[x] == "isRequired" {
					communicationFields = append(communicationFields, comm.Fields(prefix, "validate_required"))
				}
			}
		} else {
			for i := 0; i < quantityElements; i++ {
				prefixField := fmt.Sprintf("%s[%#v]", prefix, i)
				valueField := elementsOfSlice.Index(i)
				if valueField.Kind() == reflect.Struct {
					validateErrors := validateStruct(valueField, prefixField)
					communicationFields = append(communicationFields, validateErrors...)
				} else {
					validateErrors := validateField(prefixField, tagValidate, valueField)
					communicationFields = append(communicationFields, validateErrors...)
				}
			}
		}
	}
	return
}

func validateDate(prefix string, tagValidate string, valueField reflect.Value) (communicationFields []communication.Fields) {
	comm := communication.New()
	tags := getTags(tagValidate)
	if len(tags) >= 1 {
		//All values to string
		value := fmt.Sprint(valueField)
		//Efetuando looping nas tags e efetuando as validações
		for x := range tags {
			if tags[x] == "isRequired" {
				valid, _ := validator.IsDate(value)
				fmt.Println(validator.IsDate(value))
				if !valid {
					communicationFields = append(communicationFields, comm.Fields(prefix, "validate_date"))
				}
			}
		}
	}
	return
}

func validateField(prefix string, tagValidate string, valueField reflect.Value) (communicationFields []communication.Fields) {
	comm := communication.New()
	tags := getTags(tagValidate)
	if len(tags) >= 1 {
		//All values to string
		value := fmt.Sprint(valueField)

		//Efetuando looping nas tags e efetuando as validações
		for x := range tags {
			if tags[x] == "isPassword" && value != "" {
				check := govalidator.IsByteLength(value, 6, 40)
				if !check {
					communicationFields = append(communicationFields, comm.Fields(prefix, "validate_password_length"))
				}
			} else if tags[x] == "isNotZero" {
				check := value != "0" && value != "0.0"
				if !check {
					communicationFields = append(communicationFields, comm.Fields(prefix, "validate_required"))
				}
			} else {

				simpleValidations, exists := simpleValidator[tags[x]]

				if exists {
					check := simpleValidations.validator(value)
					if !check {
						communicationFields = append(communicationFields, communication.Fields{Field: prefix, Message: simpleValidations.message, Code: simpleValidations.code})
					}
				}
			}
		}
	}
	return
}

func getTags(tagValidate string) (tags []string) {
	if tagValidate != "" {
		//Broken tags per |
		tags = strings.Split(tagValidate, "|")

		//If not broken per | so try per ;
		if len(tags) == 1 {
			tags = strings.Split(tags[0], ";")
		}
	}
	return
}
