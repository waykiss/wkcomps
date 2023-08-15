package number

import (
	"fmt"
	"github.com/waykiss/wkcomps/list"
	"strconv"
	"strings"
)

var ext = []map[int]string{
	{
		1: "um", 2: "dois", 3: "três", 4: "quatro", 5: "cinco", 6: "seis", 7: "sete", 8: "oito", 9: "nove", 10: "dez",
		11: "onze", 12: "doze", 13: "treze", 14: "quatorze", 15: "quinze",
		16: "dezesseis", 17: "dezessete", 18: "dezoito", 19: "dezenove",
	},
	{
		2: "vinte", 3: "trinta", 4: "quarenta", 5: "cinquenta", 6: "sessenta", 7: "setenta", 8: "oitenta",
		9: "noventa",
	},
	{
		1: "cento", 2: "duzentos", 3: "trezentos", 4: "quatrocentos", 5: "quinhentos", 6: "seissentos",
		7: "setessentos", 8: "oitocentos", 9: "novecentos",
	},
}

var und = []string{"", " mil", " milhão", " milhões", " bilhão", " bilhões", " trilhão", " trilhões"}

type extensoUtil struct{}

var extenso = extensoUtil{}

func (e extensoUtil) Get(intValueStr, floatValueStr string) string {
	//feita a validacao sobre inteiros dessa forma pois o IsInt quando recebe "00" retorna false
	intValueStr = strings.ReplaceAll(intValueStr, ".", "")
	_, err1 := strconv.Atoi(intValueStr)
	_, err2 := strconv.Atoi(floatValueStr)
	if err1 != nil || err2 != nil {
		return ""
	}
	var ret []string
	grand := 0
	floatIntValue, _ := strconv.Atoi(floatValueStr)
	if floatIntValue == 0 {
		ret = append(ret, "zero centavos")
	} else if floatIntValue == 1 {
		ret = append(ret, "um centavo")
	} else {
		ret = append(ret, e.getCent(floatValueStr, 0)+" centavos")
	}
	intValueInt, _ := strconv.Atoi(intValueStr)
	if intValueInt == 0 {
		ret = append(ret, "zero reais")
		ret = list.ArrayStringReverse(ret)
		return fmt.Sprintf("%s e %s", ret[0], ret[1])
	} else if intValueInt == 1 {
		ret = append(ret, "um real")
		ret = list.ArrayStringReverse(ret)
		return fmt.Sprintf("%s e %s", ret[0], ret[1])
	}

	for intValueStr != "" {
		s := ""
		if len(intValueStr) < 3 {
			s = intValueStr
			intValueStr = ""
		}
		if len(intValueStr) >= 3 {
			s = intValueStr[len(intValueStr)-3:]
			intValueStr = intValueStr[:len(intValueStr)-3]
		}
		if grand == 0 {
			ret = append(ret, e.getCent(s, grand)+" reais")
		} else {
			ret = append(ret, e.getCent(s, grand))
		}
		grand += 1
		if intValueStr == "" && len(ret) > 1 {
			if len(ret[grand]) > 2 && ret[grand] == "um mil" {
				ret[grand] = "mil"
			}
		}
	}
	ret = list.ArrayStringReverse(ret)
	r := ""
	for i := range ret {
		v := ret[i]
		valorInterno := ""
		switch i {
		case 0:
			r = v
		case len(ret) - 1:
			r = fmt.Sprintf("%s e %s", r, v)
		default:
			valorInterno += ret[i]
			if len(ret) > 3 && i == (len(ret)-2) {
				if valorInterno == " reais" {
					r = fmt.Sprintf("%s %s", r, "de")
				} else if v == "cem reais" {
					r = fmt.Sprintf("%s %s", r, "e")
				}
			}
			r = fmt.Sprintf("%s %s", r, v)

		}

	}
	return r
}

func (e extensoUtil) getCent(s string, grand int) string {
	normalize := func(r, s string, grand int) (ret string) {
		if s == "100" {
			r = "cem"
		}
		switch grand {
		case 1:
		case 2:
			if r != "um" {
				grand++
			}
		case 3:
			grand = 4
			if r != "um" {
				grand++
			}
		case 4:
			grand = 6
			if r != "um" {
				grand++
			}
		default:
			ret = r
			return
		}
		ret = fmt.Sprintf("%s%s", r, und[grand])

		return
	}

	aux := s
	zero := "0"
	result := ""
	for i := 0; i < (3 - len(s)); i++ {
		result = fmt.Sprintf("%s%s", result, zero)
	}
	s = fmt.Sprintf("%s%s", result, aux)

	if s == "000" {
		return ""
	}
	if s == "100" && grand == 0 {
		return "cem"
	}
	ret := ""
	dez := fmt.Sprintf("%c%c", s[1], s[2])
	if fmt.Sprintf("%c", s[0]) != "0" {
		aux, _ := strconv.Atoi(fmt.Sprintf("%c", s[0]))
		ret = fmt.Sprintf("%s%s", ret, ext[2][aux])
		if fmt.Sprintf("%s", dez) != "00" {
			ret += " e "
		} else {
			return normalize(ret, s, grand)
		}
	}
	dezInt, _ := strconv.Atoi(dez)
	if dezInt < 20 {
		ret += ext[0][dezInt]
	} else {
		if fmt.Sprintf("%c", s[1]) != "0" {
			position1, _ := strconv.Atoi(fmt.Sprintf("%c", s[1]))
			ret = fmt.Sprintf("%s%s", ret, ext[1][position1])
			if fmt.Sprintf("%c", s[2]) != "0" {
				position2, _ := strconv.Atoi(fmt.Sprintf("%c", s[2]))
				ret = fmt.Sprintf("%s e %s", ret, ext[0][position2])
			}
		}
	}

	return normalize(ret, s, grand)
}
