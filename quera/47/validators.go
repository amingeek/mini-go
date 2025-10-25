package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func intToString(i int) string {
	return strconv.Itoa(i)
}

func list(n int) []int {
	res := []int{}
	for _, ch := range intToString(n) {
		res = append(res, int(ch-'0'))
	}
	return res
}

func sum(listn []int) int {
	sums := 0
	for _, ch := range listn {
		sums += ch
	}
	return sums
}

var validUsernameLen validator.Func = func(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return len(data) >= 8 && len(data) <= 32
}

var validateLenMin validator.Func = func(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return len(data) >= 2 && len(data) <= 16
}

var validatePhoneNumberLen validator.Func = func(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return len(data) == 11
}

var validBirthDate validator.Func = func(fl validator.FieldLevel) bool {
	birthDate, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	now := time.Now()
	sevenYearsAgo := now.AddDate(-7, 0, 0)
	hundredYearsAgo := now.AddDate(-100, 0, 0)

	if birthDate.After(sevenYearsAgo) {
		return false
	}
	if birthDate.Before(hundredYearsAgo) {
		return false
	}

	return true
}

var validNationalID validator.Func = func(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if len(data) != 10 {
		return false
	}

	numbers := []int{}
	for _, ch := range data {
		if ch < '0' || ch > '9' {
			return false
		}
		numbers = append(numbers, int(ch-'0'))
	}

	weights := []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 9; i++ {
		sum += numbers[i] * weights[i]
	}

	remainder := sum % 11
	control := numbers[9]

	if remainder < 2 {
		return control == remainder
	}
	return control == (11 - remainder)
}

var validateLenEmail validator.Func = func(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return len(data) >= 8 && len(data) <= 32
}
var startsWith09 validator.Func = func(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return len(data) == 11 && data[:2] == "09"
}

func registerValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("lendata", validateLenMin)
		v.RegisterValidation("nationalid", validNationalID)
		v.RegisterValidation("birthdate", validBirthDate)
		v.RegisterValidation("lenusername", validUsernameLen)
		v.RegisterValidation("phone_number", validatePhoneNumberLen)
		v.RegisterValidation("lenemail", validateLenEmail)
		v.RegisterValidation("startswith09", startsWith09)
	}
}

func init() {
	registerValidators()
}
