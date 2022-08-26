package validation

import (
	"github.com/Nhanderu/brdoc"
	"github.com/badoux/checkmail"
	"github.com/rodrigorodriguescosta/govalidator"
)

//IsCpfCnpjValid verifica se o valor passado é um cpf ou cnpj valido
func IsCpfCnpjValid(value string) bool {
	valid := brdoc.IsCPF(value)
	if !valid {
		valid = brdoc.IsCNPJ(value)
	}
	return valid
}

func IsCpfValid(value string) bool {
	return brdoc.IsCPF(value)
}

func IsCnpjValid(value string) bool {
	return brdoc.IsCNPJ(value)
}

func IsEmailValid(value string) bool {
	err := checkmail.ValidateFormat(value)
	return err == nil
}

func IsOnlyNumber(str string) bool {
	return govalidator.IsNumeric(str)
}

func IsIn(str string, params ...string) bool {
	for _, param := range params {
		if str == param {
			return true
		}
	}
	return false
}

//IsInInt verifica se um valor inteiro contem no array/slice passado como parametro
func IsInInt(str int, params ...int) bool {
	for _, param := range params {
		if str == param {
			return true
		}
	}
	return false
}

// IsFilled check if field has a value, if so, check the length
func IsFilled(value string, min, max int) bool {
	if govalidator.IsNull(value) {
		return false
	}
	return IsByteLength(value, min, max)
}

func IsValidId(value string) bool {
	return IsFilled(value, 26, 37)
}

//IsByteLength check length of the string, if the string is empty, skip the validation
func IsByteLength(str string, min, max int) bool {
	if str != "" {
		return govalidator.IsByteLength(str, min, max)
	}
	return true
}
