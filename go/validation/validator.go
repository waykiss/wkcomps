package validation

import (
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/rodrigorodriguescosta/govalidator"
	"golang.org/x/text/language"
	"strconv"
	"strings"
	"time"
)

// all constants for all validations
const (
	isByteLength  = "isByteLength"
	inRangeFloat  = "inRangeFloat"
	inRangeInt    = "inRangeInt"
	equalsFloat   = "equalsFloat"
	isRequired    = "isRequired"
	isObjectId    = "isObjectId"
	isEmailValid  = "isEmailValid"
	isOnlyNumber  = "isOnlyNumber"
	isGreaterThan = "isGreaterThan"
	isNotValid    = "isNotValid"
	isNotUrlValid = "isNotUrlValid"
)

// Constants for time format
const (
	DT14               = "20060102150405"
	DT8                = "20060102"
	DT6                = "200601"
	MonthDay           = "1/2"
	RFC3339FullDate    = "2006-01-02"
	RFC3339Milli       = "2006-01-02T15:04:05.999Z07:00"
	ISO8601            = "2006-01-02T15:04:05Z0700"
	ISO8601TZHour      = "2006-01-02T15:04:05Z07"
	ISO8601NoTZ        = "2006-01-02T15:04:05"
	ISO8601MilliNoTZ   = "2006-01-02T15:04:05.000"
	ISO8601CompactZ    = "20060102T150405Z0700"
	ISO8601CompactNoTZ = "20060102T150405"
	ISO8601YM          = "2006-01"
	// MySQL, BigQuery, etc.
	SQLTimestamp     = "2006-01-02 15:04:05"
	SQLTimestampPgTz = "2006-01-02 15:04:05.999999-07"
	// GMT time in format dd:mm:yy hh:mm
	DateDMYHM2 = "02:01:06 15:04"
	BrFormat1  = "02/01/2006"
	BrFormat2  = "02/01/2006 15:04"
	BrFormat3  = "02/01/2006 15:04:05"
	UsFormat1  = "2006/01/02"
	UsFormat2  = "2006/01/02 15:04"
	UsFormat3  = "2006/01/02 15:04:05"
)

// Var custom time format
var datetimeCustomFormats = [22]string{
	BrFormat1, BrFormat2, BrFormat3, UsFormat1, UsFormat2, UsFormat3,
	DT14, DT8, DT6, MonthDay, RFC3339FullDate, RFC3339Milli, ISO8601, ISO8601TZHour, ISO8601NoTZ, ISO8601MilliNoTZ,
	ISO8601CompactZ, ISO8601CompactNoTZ, ISO8601YM, SQLTimestamp, SQLTimestampPgTz, DateDMYHM2,
}

//struct to concat errors validation
type Validation struct {
	errors   []string
	language language.Tag
}

func NewValidation(language language.Tag) *Validation {
	return &Validation{
		language: language,
	}
}

// Error return the error instance
func (v *Validation) Error() (err error) {
	var sb strings.Builder
	if len(v.errors) == 0 {
		return
	}
	for _, msg := range v.errors {
		sb.WriteString(fmt.Sprintf("%s,", msg))
	}
	r := strings.TrimSuffix(sb.String(), ",")
	return fmt.Errorf("%s", r)
}

//GetErrors retorna os erros de validacao, após executar essa funcao, os erros são zerados
func (v *Validation) GetErrors() []string {
	defer func(v *Validation) {
		v.errors = nil
	}(v)
	return v.errors
}

//Validated funcao para retornar se existem erros ou nao
func (v *Validation) Validated() bool {
	return len(v.errors) == 0
}

func (v *Validation) AddMessage(str string, params ...interface{}) {
	v.errors = append(v.errors, fmt.Sprintf(str, params...))
}

//AddErroMessage essa funcao cria uma mensagem de erro baseado no objeto erro
func (v *Validation) AddErroMessage() {
	v.errors = append(v.errors, "Ocorreu um erro no sistema e nossa equipe já foi notificada, aguarde que já estamos "+
		"ciente e já trabalhando na solução.")
}

//AddFirstMessage adiciona uma mensagem para a primeira na lista do slice(array)
func (v *Validation) AddFirstMessage(str string, params ...interface{}) {
	newSlice := []string{fmt.Sprintf(str, params...)}
	v.errors = append(newSlice, v.errors...)
}

//SetLastMessage altera a ultima mensagem de validacao
func (v *Validation) SetLastMessage(msg string) {
	v.errors[len(v.errors)-1] = msg
}

func (v *Validation) IsIn(msg string, str string, params ...string) bool {
	if IsIn(str, params...) {
		return true
	}

	// verifique se tem o nome `precisa ser` na string de mensagem, se tiver, então nao use a mensagem padrao
	if !strings.Contains(msg, "precisa ser") {
		msg = fmt.Sprintf("%s precisa ser uma das seguintes opções : %s, valor recebido foi '%s'", msg,
			strings.Join(params, ","),
			str)
	}
	v.errors = append(v.errors, msg)
	return false
}

//IsByteLength check length of the string, if the string is empty, skip the validation
func (v *Validation) IsByteLength(msg, str string, min, max int) bool {
	if str != "" {
		if !IsByteLength(str, min, max) {
			msg = fmt.Sprintf(v.getMessage(isByteLength), msg, str, min, max, len(str))
			v.errors = append(v.errors, msg)
			return false
		}
	}
	return true
}

//EqualsFloat verifica se os valores recebidos são iguais
func (v *Validation) EqualsFloat(msg string, value, equalValue float64) bool {
	if value != equalValue {
		msg = fmt.Sprintf(v.getMessage(equalsFloat), msg, equalValue, value)
		v.errors = append(v.errors, msg)
		return false
	}
	return true
}

func (v *Validation) InRangeFloat32(msg string, value, min, max float32) bool {
	if !govalidator.InRangeFloat32(value, min, max) {
		msg = fmt.Sprintf(v.getMessage(inRangeFloat), msg, min, max, value)
		v.errors = append(v.errors, msg)
		return false
	}
	return true
}

func (v *Validation) InRangeFloat64(msg string, value, min, max float64) bool {
	if !govalidator.InRangeFloat64(value, min, max) {
		msg = fmt.Sprintf(v.getMessage(inRangeFloat), msg, min, max, value)
		v.errors = append(v.errors, msg)
		return false
	}
	return true
}

func (v *Validation) InRangeInt(msg string, value, min, max int) bool {
	if !govalidator.InRangeInt(value, min, max) {
		msg = fmt.Sprintf(v.getMessage(inRangeInt), msg, min, max, value)
		v.errors = append(v.errors, msg)
		return false
	}
	return true
}

// isGreaterThanFloat64 check if int value is greater than
func (v *Validation) IsGTFloat64(msg string, value, min float64) bool {
	if value < min {
		msg = fmt.Sprintf(v.getMessage(isGreaterThan), msg, strconv.FormatFloat(min, 'f', -1, 64))
		v.errors = append(v.errors, msg)
		return false
	}
	return true
}

// isGreaterThanInt check if int value is greater than
func (v *Validation) IsGTInt(msg string, value, min int) bool {
	if value < min {
		msg = fmt.Sprintf(v.getMessage(isGreaterThan), msg, strconv.Itoa(value))
		v.errors = append(v.errors, msg)
		return false
	}
	return true
}

// isGreaterThanInt check if int value is greater than
func (v *Validation) IsGTTime(msg string, value, min time.Time) bool {
	if value.Before(min) {
		msg = fmt.Sprintf(v.getMessage(isGreaterThan), msg, min, value)
		v.errors = append(v.errors, msg)
		return false
	}
	return true
}

func (v *Validation) IsObjectId(fieldName string, value string) bool {
	if !govalidator.IsMongoID(value) {
		msg := fmt.Sprintf(v.getMessage(isObjectId), fieldName)
		v.errors = append(v.errors, msg)
		return false
	}
	return true
}

// IsFilled check if field has a value, if so, check the length
func (v *Validation) IsFilled(fieldName string, value string, min, max int) bool {
	if govalidator.IsNull(value) {
		msg := fmt.Sprintf(v.getMessage(isRequired), fieldName)
		v.errors = append(v.errors, msg)
		return false
	}
	return v.IsByteLength(fieldName, value, min, max)
}

// IsFilledTime verifica se um valor do tipo time foi preenchido
func (v *Validation) IsFilledTime(fieldName string, value time.Time) bool {
	if value.IsZero() {
		msg := fmt.Sprintf(v.getMessage(isRequired), fieldName)
		v.errors = append(v.errors, msg)
		return false
	}
	return true
}

func (v *Validation) IsValidId(field, value string) bool {
	return v.IsFilled(field, value, 26, 37)
}

func (v *Validation) IsObjectIdAndFilled(fieldName string, value string) bool {
	if ok := v.IsFilled(fieldName, value, 24, 24); ok == true {
		return v.IsObjectId(fieldName, value)
	}
	return false
}

func (v *Validation) IsValidEmailFormat(fieldName string, value string) bool {
	if !IsEmailValid(value) {
		msg := fmt.Sprintf(v.getMessage(isEmailValid), fieldName, value)
		v.errors = append(v.errors, msg)
	}
	return true
}

//IsOnlyNumber returns true if the string contains only number
func (v *Validation) IsOnlyNumber(fieldName, str string) bool {
	onlyNumber := IsOnlyNumber(str)
	if !onlyNumber {
		msg := fmt.Sprintf(v.getMessage(isOnlyNumber), fieldName)
		v.errors = append(v.errors, msg)
	}
	return onlyNumber
}

//IsCpfCnpjValid check if Cpf or Identification are valid
func (v *Validation) IsCpfCnpjValid(fieldName, value string) bool {
	valid := IsCpfValid(value)
	if !valid {
		valid = IsCnpjValid(value)
	}
	if !valid {
		msg := fmt.Sprintf(v.getMessage(isNotValid), fieldName, value)
		v.errors = append(v.errors, msg)
	}
	return valid
}

func (v *Validation) IsCpfValid(fieldName, value string) bool {
	valid := IsCpfValid(value)
	if !valid {
		msg := fmt.Sprintf(v.getMessage(isNotValid), fieldName, value)
		v.errors = append(v.errors, msg)
	}
	return valid
}

func (v *Validation) IsCnpjValid(fieldName, value string) bool {
	valid := IsCnpjValid(value)
	if !valid {
		msg := fmt.Sprintf(v.getMessage(isNotValid), fieldName, value)
		v.errors = append(v.errors, msg)
	}
	return valid
}

func (v *Validation) getMessage(messageId string) string {
	switch lang := v.language; lang {
	default:
		return messagesPtBR[messageId]
	}
}

//IsInt Verifica se a string passada é um numero inteiro
func (v *Validation) IsInt(filter string) bool {
	if _, err := strconv.Atoi(filter); err == nil {
		return true
	}
	return false
}

//IsDateTime check if a datetime in string is valid
func (v *Validation) IsDateTime(datetime string) bool {
	for _, format := range datetimeCustomFormats {
		if _, err := time.Parse(format, datetime); err == nil {
			return true
		}
	}
	return false
}

//IsStateBR check if a state is valid
func (v *Validation) IsStateBR(fieldName, value string) bool {
	isValid := IsStateBR(value)
	if !isValid {
		v.errors = append(v.errors, fmt.Sprintf(v.getMessage(isNotValid), fieldName, value))
	}
	return isValid
}

//IsUrl check if the string is a valida url
func (v *Validation) IsUrl(fieldName, value string) bool {
	isValid := govalidator.IsURL(value)
	if !isValid {
		v.errors = append(v.errors, fmt.Sprintf(v.getMessage(isNotUrlValid), fieldName, value))
	}
	return isValid
}

//IsUrl check if the string is a valida url
func (v *Validation) IsUUID(fieldName, value string) bool {
	isValid := govalidator.IsUUID(value)
	if !isValid {
		v.errors = append(v.errors, fmt.Sprintf(v.getMessage(isNotValid), fieldName, value))
	}
	return isValid
}

// IsCreditCardNumber valida se um número de cartao é valido
func (v *Validation) IsCreditCardNumber(fieldName, number string) bool {

	var sum int
	var alternate bool

	// Gets the Card number length
	numberLen := len(number)

	// For numbers that is lower than 13 and
	// bigger than 19, must return as false
	if numberLen < 13 || numberLen > 19 {
		return false
	}

	// Parse all numbers of the card into a for loop
	for i := numberLen - 1; i > -1; i-- {
		// Takes the mod, converting the current number in integer
		mod, _ := strconv.Atoi(string(number[i]))
		if alternate {
			mod *= 2
			if mod > 9 {
				mod = (mod % 10) + 1
			}
		}

		alternate = !alternate
		sum += mod
	}

	isValid := sum%10 == 0
	if !isValid {
		v.errors = append(v.errors, fmt.Sprintf(v.getMessage(isNotValid), fieldName, number))
	}

	return isValid
}

//IsEmail verifique se uma string é um email válido
func IsEmail(v string) bool {
	if err := checkmail.ValidateFormat(v); err != nil {
		return false
	}
	return true
}
