package str

import (
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/andreyvit/diff"
	"github.com/iancoleman/strcase"
	"github.com/rodrigorodriguescosta/govalidator"
	"github.com/segmentio/ksuid"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// StringInSlice check if certain string contain in slice
func StringInSlice(a string, list ...string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type RandStringChars string

const (
	FormatCpf        string = "###.###.###-##"
	FormatCnpj       string = "##.###.###/####-##"
	FormatCreditCard string = "####.####.####.####"
	MaskCreditCard   string = "****.****.****-####"
)

func (r RandStringChars) String() string {
	return string(r)
}

const RandStringCharsOnlyNumbers RandStringChars = "0123456789"
const RandStringCharsOnlyLetters RandStringChars = "abcdefghijklmnopqrstuvwxyz"
const RandStringCharsLettersAndNumbers RandStringChars = "abcdefghijklmnopqrstuvwxyz0123456789"
const RandStringCharsAllLettersAndNumbers RandStringChars = `0123456789ABCDEFGHIJKLMNOPQRSTUVW` +
	`XYZabcdefghijklmnopqrstuvwxyz`

const (
	RegexOnlyLetters = "[^A-Z]"
	RegexOnlyNumbers = "[^0-9]"
)

// RandString generate a random string
// Use the constants RandStringCharsOnlyNumbers, RandStringCharsOnlyLetters and RandStringCharsLettersAndNumbers
// to facilitate to pass the type of character you'd like to generate
// limit the number of characters at 250
func RandString(n int, chars string) (r string) {
	if n <= 0 {
		return
	}
	if n > 250 {
		n = 250
	}
	result := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// RemoveExtraSpaces remove all duplicated spaces over whole string
func RemoveExtraSpaces(v string) string {
	return strings.Join(strings.Fields(v), " ")
}

// RemoveAccent remove all accents from given string
func RemoveAccent(v string) string {
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, v)
	return result
}

// NoSpaceNoAccent remove all the spaces of the given string, begin, end and middle of the string
func NoSpaceNoAccent(v string) string {
	v = RemoveExtraSpaces(v)
	v = RemoveAccent(v)
	return v
}

// TrimAll remove todos os espacos em uma string
func TrimAll(v string) string {
	return strings.ReplaceAll(v, " ", "")
}

// UpperNoSpaceNoAccent apply 3 steps over string, transform to uppercase, remove duplicate spaces
// and remove accents
func UpperNoSpaceNoAccent(v string) string {
	v = strings.ToUpper(v)
	v = NoSpaceNoAccent(v)
	return v
}

// LowerNoSpaceNoAccent apply 3 steps over string, transform to lowercase, remove duplicate spaces
// and remove accents
func LowerNoSpaceNoAccent(v string) string {
	v = RemoveExtraSpaces(v)
	v = strings.ToLower(v)
	v = RemoveAccent(v)
	return v
}

/*
	ReplaceSpecialCharacters substitui todos os caracteres especiais do primeiro parametro pelo valor informado no

segundo parametro.
Se o segundo parametro nao for informado, os caracteres serão apenas removidos.
A partir do terceiro parametro sao consideradas excessoes a regra caso necessario.
Ex:
ReplaceSpecialCharacters("_teste!")
saida: teste
ReplaceSpecialCharacters("_teste!", "@")
saida: @teste@
ReplaceSpecialCharacters("_teste!", "@", "!")
saida: @teste!
*/
func ReplaceSpecialCharacters(v ...string) (result string) {
	substitute := ""
	exceptions := ""
	for i := 0; i < len(v); i++ {
		switch i {
		case 0:
			result = v[i]
			continue
		case 1:
			substitute = v[i]
			continue
		default:
			exceptions += v[i]
			continue
		}
	}
	result = RemoveAccent(result)
	re := regexp.MustCompile(fmt.Sprintf(`[^A-Za-z0-9%s]`, exceptions))
	result = re.ReplaceAllString(result, substitute)

	return
}

// OnlyNumbers returns only numbers given a string
func OnlyNumbers(v string) string {
	if v != "" {
		reg, _ := regexp.Compile("[^0-9]+")
		return reg.ReplaceAllString(v, "")
	}
	return v
}

// ContainNumber verifica se uma string contem numeros
func ContainNumber(v string) bool {
	if v != "" {
		reg, _ := regexp.Compile("[0-9]")
		return reg.MatchString(v)
	}
	return false
}

// PointerToString convert string pointer to string as Zero value if pointer is nil
func PointerToString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func ToString(str interface{}) string {
	if str == nil {
		return ""
	}
	if value, ok := str.(string); ok {
		return value
	}
	return ""
}

// StringPToUint convert string pointer to uint, returns zero if pointer is nil
func StringPToUint(v *string) uint {
	if v == nil {
		return 0
	}

	u64, err := strconv.ParseUint(*v, 10, 32)
	if err != nil {
		return 0
	}
	return uint(u64)
}

// GetFilenameFromPath return filename given the path(url or local path)
func GetFilenameFromPath(p string) string {
	return path.Base(p)
}

// GetFilenameWithoutExtension returns the filename without extension of given path
func GetFilenameWithoutExtension(v string) string {
	return strings.TrimSuffix(v, path.Ext(v))
}

// GetPathFromFilename get only the path of the whole path passed by parameter, it returns path without filename
func GetPathFromFullFilename(v string) string {
	return strings.TrimSuffix(v, path.Ext(v))
}

// IsUrl check if the string is a valid url
func IsUrl(url string) bool {
	return govalidator.IsURL(url)
}

// IsInt check if a string is a valid number
func IsInt(url string) bool {
	return govalidator.IsInt(url)
}

func IsOnlyNumbers(v string) bool {
	if posV := OnlyNumbers(v); posV == v {
		return true
	}
	return false
}

// ContainNumbers verifica se a string passada contem numeros entre os caracteres
func ContainNumbers(v string) bool {
	result := RemoveAccent(v)
	re := regexp.MustCompile(RegexOnlyNumbers)
	result = re.ReplaceAllString(result, "")
	return result != ""
}

// ContainLetters verifica se a string passada contem letras entre os caracteres
func ContainLetters(v string) bool {
	result := RemoveAccent(v)
	re := regexp.MustCompile(RegexOnlyLetters)
	result = re.ReplaceAllString(result, "")
	return result != ""
}

// ContainNumbersAndLetters verifica se a string passada contem letras e numeros ao mesmo tempo
func ContainNumbersAndLetters(v string) (r bool) {
	r = ContainNumbers(v)
	if !r {
		return
	}
	r = ContainLetters(v)
	return
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToLowerCamelCase(v string) (result string) {
	var convertNext bool
	for idx, char := range v {
		if idx == 0 {
			result += strings.ToLower(string(char))
			continue
		}
		if char == '_' {
			convertNext = true
			continue
		} else if unicode.IsUpper(char) {
			if idx+1 < len(v) && unicode.IsLower(rune(v[idx+1])) {
				result += string(char)
				continue
			}
		}

		if convertNext {
			convertNext = false
			result += strings.ToUpper(string(char))
		} else {
			result += strings.ToLower(string(char))
		}
	}
	return
}

func ToCamel(v string) string {
	return strcase.ToCamel(v)
}

func Uuid() string {
	return ksuid.New().String()
}

// GetLastCharacter retorna os ultimos x/count caracteres de uma string
func GetLastCharacter(v string, count int) (r string) {
	if v != "" && count < len(v) {
		r = v[len(v)-count:]
	}

	return
}

func ExtractValue(body string, key string) string {
	keystr := "\"" + key + "\":[^,;\\]}]*"
	r, _ := regexp.Compile(keystr)
	match := r.FindString(body)
	keyValMatch := strings.Split(match, ":")
	return strings.ReplaceAll(keyValMatch[1], "\"", "")
}

// ReplaceAtIndex altera um caractere em uma string baseado no index da mesma
func ReplaceAtIndex(in string, char rune, idx int) string {
	out := []rune(in)
	out[idx] = char
	return string(out)
}

// Format essa funcao tem como objetivo mascarar strings em um determinado formato, o caractere # será substituído pelo
// valor subsequente da string, o caractere * a funcao entende como uma mascara em si, ou seja, ela conta os caracteres
func Format(v, format string) (result string) {
	// nao processe caso a quantidade de caracteres for diferente dos valores que devem ser preenchidos na mascara
	if (strings.Count(format, "#") + strings.Count(format, "*")) != len(v) {
		result = v
		return
	}
	idx := 0
	result = format
	for _, c := range format {
		// o asterisco, conta o índice pois se trata de mascarar
		if c == '*' {
			idx++
			continue
		}
		if c == '#' {
			charAt := v[idx]
			result = strings.Replace(result, "#", string(charAt), 1)
			idx++
			continue
		}
	}
	return
}

//
////GetJsonFromArray funcao que dado um array de string e um campo, retorna um json no formato de string
//func GetJsonFromArray(fieldName string, array []string) string {
//	var arrayQuery []map[string]string
//	for _, value := range array {
//		arrayQuery = append(arrayQuery, map[string]string{fieldName: value})
//	}
//	itemArrayString, _ := json.Marshal(arrayQuery)
//	return string(itemArrayString)
//}

// GetNotZeroValue dado um array de valores string, é retornado o que tem valor nao zerovalue
func GetNotZeroValue(p ...string) string {
	for _, v := range p {
		if v != "" {
			return v
		}
	}
	return ""
}

// StrPad returns the input string padded on the left, right or both sides using padType to the specified padding length padLength.
//
// Example:
// input := "Codes";
// StrPad(input, 10, " ", "RIGHT")        // produces "Codes     "
// StrPad(input, 10, "-=", "LEFT")        // produces "=-=-=Codes"
// StrPad(input, 10, "_", "BOTH")         // produces "__Codes___"
// StrPad(input, 6, "___", "RIGHT")       // produces "Codes_"
// StrPad(input, 3, "*", "RIGHT")         // produces "Codes"
func StrPad(input string, padLength int, padString string, padType string) string {
	var output string

	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	switch padType {
	case "RIGHT":
		output = input + strings.Repeat(padString, int(repeat))
		output = output[:padLength]
	case "LEFT":
		output = strings.Repeat(padString, int(repeat)) + input
		output = output[len(output)-padLength:]
	case "BOTH":
		length := (float64(padLength - inputLength)) / float64(2)
		repeat = math.Ceil(length / float64(padStringLength))
		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(float64(length)))] + input + strings.Repeat(padString, int(repeat))[:int(math.Ceil(float64(length)))]
	}

	return output
}

// StrPadRight adicionar `padLength` quantidade de caracteres a direita de uma string
func StrPadRight(input string, padLength int, padString string) string {
	return StrPad(input, padLength, padString, "RIGHT")
}

// CutString remove parte do texto de acordo com o inicio e fim da sequencia
// informados, retornando o texto atualizado
func CutString(str, initial, final string, n int) (strOut string) {
	strCut := GetStringBetween(str, initial, final, true)
	//efetua substituição da string capturada por uma string vazia
	if strCut == "" {
		return str
	}
	strOut = strings.Replace(str, strCut, "", n)
	return
}

/*
GetStringBetween retorna a string que está entre a string inicial e final passada como parametro, caso nao encontre,
retorne uma string vazia
*/
func GetStringBetween(str, initial, final string, includeInitialAndFinal bool) (strOut string) {
	iniIndex := strings.Index(str, initial)
	//se nao encontrar o index da string a ser removida retornar a propria string base
	if iniIndex < 0 {
		return
	}
	finIndex := iniIndex + strings.Index(str[iniIndex:], final)
	if finIndex <= iniIndex {
		return
	}
	//captura a sequencia a ser removida
	strOut = str[iniIndex : finIndex+len(final)]
	if includeInitialAndFinal == false {
		strOut = str[iniIndex+(len(initial)) : finIndex]
	}
	return
}

// GetDiff returns the differences between the two given strings
func GetDiff(str1, str2 string) (strOut string) {
	result := diff.CharacterDiff(str1, str2)
	return GetStringBetweenAll(result, "(", ")", true)
}

// GetStringBetweenAll Returns all strings between the start and end strings
func GetStringBetweenAll(str, initial, final string, includeInitialAndFinal bool) (strOut string) {
	//quantas vezes a sequencia inicial se repete para ser buscada no loop
	initialX := strings.Count(str, initial)
	for i := 0; i < initialX; i++ {
		iniIndex := strings.Index(str, initial)
		// //se nao encontrar o index da string a ser removida retornar a propria string base
		// if iniIndex < 0 {
		// 	return
		// }
		finIndex := iniIndex + strings.Index(str[iniIndex:], final)
		if finIndex <= iniIndex {
			return
		}
		//captura a sequencia a ser removida
		res := str[iniIndex : finIndex+len(final)]
		if includeInitialAndFinal == false {
			res = str[iniIndex+(len(initial)) : finIndex]
		}
		strOut = fmt.Sprintf("%s %s", strOut, res)
		//elimina a sequencia ja analizada
		str = str[finIndex+len(final):]
	}

	return
}

// GetUriFromUrl retorna uma URI dado uma URL
func GetUriFromUrl(urlString string) string {
	u, _ := url.Parse(urlString)
	return u.RequestURI()
}

// WordCount return the number of the words that contains in string
func WordCount(v string) int {
	words := strings.Fields(v)
	return len(words)
}

// Contains verifica se alguma das strings passada contem no primeiro parametro
func Contains(v string, strs ...string) bool {
	for _, str := range strs {
		if strings.Contains(v, str) {
			return true
		}
	}
	return false
}
