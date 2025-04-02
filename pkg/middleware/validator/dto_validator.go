
package validator

import (
	"reflect"
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
)

const MAX_NO_OF_SEATS_PER_BOOKING = 10

type DtoValidator struct {
	sync     sync.Once
	validate *validator.Validate
}

func (d *DtoValidator) Engine() interface{} {
	d.lazyInit()
	return d.validate
}

func (d *DtoValidator) ValidateStruct(any interface{}) error {
	if dataType(any) == reflect.Struct {
		d.lazyInit()
		if err := d.validate.Struct(any); err != nil {
			return error(err)
		}
	}
	return nil
}

func (d *DtoValidator) lazyInit() {
	d.sync.Do(func() {
		d.validate = validator.New()
		d.validate.SetTagName("binding")
		d.validate.RegisterValidation("phoneNumber", validatePhoneNumber())
		d.validate.RegisterValidation("maxSeats", validateMaxSeatsAllowed())
	})
}

func dataType(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	kind := value.Kind()
	if kind == reflect.Ptr {
		kind = value.Elem().Kind()
	}
	return kind
}

func validatePhoneNumber() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		phoneNumberLength := fl.Field().Len()
		if phoneNumberLength != 10 {
			return false
		}
		match, _ := regexp.MatchString("\\d{10}", fl.Field().String())
		return match
	}
}

func validateMaxSeatsAllowed() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		seatsRequested := fl.Field().Int()
		if MAX_NO_OF_SEATS_PER_BOOKING < seatsRequested {
			return false
		}
		return true
	}
}
