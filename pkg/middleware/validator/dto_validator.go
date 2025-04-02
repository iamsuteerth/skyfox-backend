package validator

import (
	"reflect"
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

		d.validate.RegisterValidation("customName", ValidateName)
		d.validate.RegisterValidation("customUsername", ValidateUsername)
		d.validate.RegisterValidation("customPhone", ValidatePhoneNumber)
		d.validate.RegisterValidation("customPassword", ValidatePassword)

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

func validateMaxSeatsAllowed() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		seatsRequested := fl.Field().Int()
		return MAX_NO_OF_SEATS_PER_BOOKING >= seatsRequested
	}
}
