package number

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/waykiss/wkcomps/str"
	"golang.org/x/exp/constraints"
	"strconv"
	"strings"
)

// Max return the maximum value given variadic argument
func Max[T constraints.Ordered](values ...T) (r T) {
	return MinMax(false, values...)
}

// Min return the maximum value given variadic argument
func Min[T constraints.Ordered](values ...T) (r T) {
	return MinMax(true, values...)
}

// MinMax return the min or max given variadic argument
func MinMax[T constraints.Ordered](min bool, values ...T) (r T) {
	if values == nil {
		return
	}
	r = values[0]
	for _, v := range values {
		if min {
			if v < r {
				r = v
			}
			continue
		}
		if v > r {
			r = v
		}
	}
	return
}

type Currency string

const (
	CurrencyBRL         Currency = "BRL"
	CurrencyBRLNoPrefix Currency = "BRL_NOPREFIX"
	CurrencyUSD         Currency = "USD"
)

// RoundFloat returns float value with precision passed by parameter
func RoundFloat(v float64, precision int) float64 {
	if precision <= 0 {
		precision = 1
	}
	format := strings.Replace("%.precisionf", "precision", strconv.Itoa(precision), precision)
	value := fmt.Sprintf(format, v)
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return val
}

// ToFloat64 converte strings and other primite types to Float64, returns 0(Zero) if
// the parameter passed is not a valid float64
func ToFloat64(i interface{}, decimals int) float64 {
	floatValue, err := strconv.ParseFloat(fmt.Sprint(i), decimals)
	if err != nil {
		return 0
	}
	return floatValue
}

// ToFloat64 converte strings and other primite types to Float64, returns 0(Zero) if
// the parameter passed is not a valid float64
func ToFloat64Pointer(i interface{}, decimals int) *float64 {
	floatValue, err := strconv.ParseFloat(fmt.Sprint(i), decimals)
	if err != nil {
		v := float64(0)
		return &v
	}
	return &floatValue
}

// PreciseNumber retorna o numero informado convertido para a precisao em casas decimais desejada
func PreciseNumber(value interface{}, precision int) (result float64) {
	result, err := strconv.ParseFloat(fmt.Sprint(value), 4)
	if err != nil {
		result = float64(0)
	}
	splitedNumber := strings.Split(fmt.Sprint(result), ".")
	if len(splitedNumber) > 1 && len(splitedNumber[1]) > precision {
		result, err = strconv.ParseFloat(splitedNumber[0]+"."+splitedNumber[1][:precision], precision)
		if err != nil {
			result = float64(0)
		}
	}
	return result
}

// Percentage retorna o resultado do calculo percentual com base no percentual, valor e precisao informados
func Percentage(value, percentage interface{}, precision int) float64 {
	var floatValue, percentageValue float64
	floatValue, err := strconv.ParseFloat(fmt.Sprint(value), 4)
	if err != nil {
		floatValue = float64(0)
	}
	percentageValue, err = strconv.ParseFloat(fmt.Sprint(percentage), 4)
	if err != nil {
		percentageValue = float64(0)
	}

	//efetuando calculo preciso da porcentagem. Necessario pois apenas usando float nem sempre gera resultado correto
	result, _ := decimal.NewFromFloat(floatValue / 100).Mul(decimal.NewFromFloat(percentageValue)).Float64()
	return RoundFloat(result, precision)
}

// CalcPercentage calcula o percentual do primeiro parametro em relacao ao segundo(total)
func CalcPercentage(total, value interface{}, precision int) float64 {
	var floatValue, totalValue float64
	floatValue, err := strconv.ParseFloat(fmt.Sprint(value), 4)
	if err != nil {
		floatValue = float64(0)
	}
	totalValue, err = strconv.ParseFloat(fmt.Sprint(total), 4)
	if err != nil {
		totalValue = float64(0)
	}
	if totalValue == 0 {
		return totalValue
	}
	return PreciseNumber((floatValue/totalValue)*100, precision)
}

// ToInt convert interface to int, if val is not a valid int value, zero will be returned
func ToInt(val interface{}) int {
	strValue := fmt.Sprint(val)
	return int(StringToFloat(strValue))
}

// StringToInt convert string type into int type, return zero value if conversion has a error
func StringToInt(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

// StringToFloat Converte uma string suspostamente com valor de float para o tipo float
func StringToFloat(val string) float64 {
	val = strings.ReplaceAll(val, ",", ".")
	i, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0
	}
	return i
}

// ToInt convert interface to int32, if val is not a valid int32 value, zero will be returned
func ToInt32(val interface{}) int32 {
	intValue, ok := val.(int32)
	if !ok {
		return 0
	}
	return intValue
}

func ToInt64(val interface{}) int64 {
	intValue, ok := val.(int64)
	if !ok {
		return 0
	}
	return intValue
}

func ToFloat32(i interface{}, decimals int) float32 {
	return float32(ToFloat64(i, decimals))
}

func PointerToFloat32(v *float32) float32 {
	if v == nil {
		return 0
	}
	return *v
}

func PointerToFloat64(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func PointerToFInt32(v *int32) int32 {
	if v == nil {
		return 0
	}
	return *v
}

func PointerToFInt16(v *int16) int16 {
	if v == nil {
		return 0
	}
	return *v
}

func PointerToInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func ToUint16(v interface{}) uint16 {
	strValue := v.(string)
	value, err := strconv.ParseUint(strValue, 10, 16)
	if err != nil {
		return 0
	}
	return uint16(value)
}

// GetFormattedMoneyBRL retorna um float formatado em dinheiro brasileiro
func GetFormattedMoneyBRL(v float64, precision int) string {
	return getFormatted(v, CurrencyBRL, precision)
}

// GetFormattedMoneyBRLWithNoPrefix retorna um float formatado em dinheiro brasileiro sem o prefixo R$
func GetFormattedMoneyBRLWithNoPrefix(v float64, precision int) string {
	return getFormatted(v, CurrencyBRLNoPrefix, precision)
}

// GetValorMonetarioPorExtenso retorna uma string contendo o significado monetario por extenso do valor informado
func GetValorMonetarioPorExtenso(v float64) string {
	intValue := AddThousandToInt(int(v))
	floatValue := GetDecimalPart(v)
	floatValue = str.StrPadRight(floatValue, 2, "0")
	if len(floatValue) > 2 {
		floatValue = floatValue[:2]
	}

	return str.RemoveExtraSpaces(extenso.Get(intValue, floatValue))
}

func GetFormattedMoneyUSD(v float64, precision int) string {
	return getFormatted(v, CurrencyUSD, precision)
}

// GetFormattedFloat retorma um float formatado
func GetFormattedFloat(v float64) string {
	intValue := AddThousandToInt(int(v))
	decimalPart := GetDecimalPart(v)
	i, _ := strconv.Atoi(decimalPart)
	if i == 0 {
		return fmt.Sprintf("%s", intValue)
	}
	return fmt.Sprintf("%s,%s", intValue, decimalPart)
}

// String returns a formatted Currency value
func getFormatted(value float64, currency Currency, precision int) string {
	intValue := AddThousandToInt(int(value))
	floatValue := GetDecimalPart(value)
	floatValue = str.StrPadRight(floatValue, 2, "0")
	if len(floatValue) > precision {
		floatValue = floatValue[:precision]
	}
	switch currency {
	case CurrencyUSD:
		intValue = strings.ReplaceAll(intValue, ".", ",")
		return fmt.Sprintf("$ %s.%s", intValue, floatValue)
	case CurrencyBRLNoPrefix:
		return fmt.Sprintf("%s,%s", intValue, floatValue)
	default:
		return fmt.Sprintf("R$ %s,%s", intValue, floatValue)
	}
}

// addThousandToInt add the Thousand separator to int value and returns a string with Thousand added
func AddThousandToInt(number int) string {
	intergerPart := fmt.Sprintf("%d", number)
	newString := ""
	aux := 0
	for i := len(intergerPart) - 1; i >= 0; i-- {
		if aux > 0 && aux%3 == 0 {
			newString = "." + newString
		}
		aux++
		newString = string(intergerPart[i]) + newString
	}
	return newString
}

// GetDecimalPart retorna a parte decimal de um float, retornando em formato string
func GetDecimalPart(value float64) (r string) {
	//extraindo a parte decimal do valor
	v := strconv.FormatFloat(value, 'f', -1, 64)
	if strings.Contains(v, ".") {
		r = strings.Split(v, ".")[1]
	}
	return
}

// GetNotZeroFloatValue dado um array de valores float64, é retornado o que tem valor maior que zero
func GetNotZeroFloatValue(p ...float64) float64 {
	for _, v := range p {
		if v > 0 {
			return v
		}
	}
	return 0.0
}

// IsIntegral Verifica se o valor informado é fracionado ou não
func IsIntegral(value float64) bool {
	return PreciseNumber(value, 0) == value
}

// FloatToIntMoney converte um float em inteiro com 2 casas decimais
func FloatToIntMoney(v float64) (r int) {
	s := strings.ReplaceAll(fmt.Sprintf("%.2f", v), ".", "")
	r, _ = strconv.Atoi(s)
	return
}
