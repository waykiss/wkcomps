package str

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestUpperNoSpaceNoAccent(t *testing.T) {
	values := []struct {
		input    string
		expected string
	}{
		{" Hey Holla ÀáÀ  ê ã", "HEY HOLLA AAA E A"},
		{"Test removing accent Á á ẽ Ê", "TEST REMOVING ACCENT A A E E"},
		{"Test removing spaces in the middle of the    words", "TEST REMOVING SPACES IN THE MIDDLE OF " +
			"THE WORDS"},
		{" Removing spaces in the edges ", "REMOVING SPACES IN THE EDGES"},
	}
	for _, v := range values {
		got := UpperNoSpaceNoAccent(v.input)
		assert.Equal(t, v.expected, got, "Conversion of the value '%f' is expected '%f' but got '%f'",
			v.input, v.expected, got)
	}
}

func TestNoSpace(t *testing.T) {
	values := []struct {
		input    string
		expected string
	}{
		{" Hey Holla ÀáÀ  ê ã", "Hey Holla ÀáÀ ê ã"},
		{" ", ""},
		{"HELLO    123", "HELLO 123"},
		{"a b C D e  F   G ", "a b C D e F G"},
		{" a ", "a"},
		{" TEST1 TESTE2  Test3 ", "TEST1 TESTE2 Test3"},
	}
	for _, v := range values {
		got := RemoveExtraSpaces(v.input)
		assert.Equal(t, v.expected, got, v.input, v.expected, got)
	}
}

func TestNoAccent(t *testing.T) {
	values := []struct {
		input    string
		expected string
	}{
		{" Hey Holla ÀáÀ  ê ã", " Hey Holla AaA  e a"},
		{" ", " "},
		{"âÂéÈÍíòÓúÚ  ", "aAeEIioOuU  "},
		{"aaa bbb  CCCC  ", "aaa bbb  CCCC  "},
	}
	for _, v := range values {
		got := RemoveAccent(v.input)
		assert.Equal(t, v.expected, got, v.input, v.expected, got)
	}
}

func TestReplaceSpecialCharacters(t *testing.T) {
	values := []struct {
		input      string
		substitute string
		exceptions []string
		expected   string
	}{
		{input: "@teste@", expected: "teste"},
		{input: "@_teste !", substitute: "_", expected: "__teste__"},
		{input: "@teste!", substitute: "_", expected: "_teste_"},
		{input: "_teste!", substitute: "[underline]", expected: "[underline]teste[underline]"},
		{input: "_teste!", substitute: "[underline]", exceptions: []string{"!"}, expected: "[underline]teste!"},
		{input: "-áAb-axxb-_@_(_%_&__$ #_`_!a!90${1}", substitute: "_", exceptions: []string{"!"}, expected: "_aAb_axxb_________________!a!90__1_"},
	}
	for _, v := range values {
		values := []string{}
		values = append(values, v.input)
		values = append(values, v.substitute)
		values = append(values, v.exceptions...)
		result := ReplaceSpecialCharacters(values...)
		assert.Equal(t, v.expected, result, "Conversion of the value '%s' is expected '%s' but got '%s'",
			v.input, v.expected, result)
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"already_snake", "already_snake"},
		{"A", "a"},
		{"AA", "aa"},
		{"AaAa", "aa_aa"},
		{"HTTPRequest", "http_request"},
		{"BatteryLifeValue", "battery_life_value"},
		{"Id0Value", "id0_value"},
		{"ID0Value", "id0_value"},
	}
	for _, test := range tests {
		have := ToSnakeCase(test.input)
		if have != test.want {
			t.Errorf("input=%q:\nhave: %q\nwant: %q", test.input, have, test.want)
		}
	}
}

func TestToLowerCamelCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"already_snake", "alreadySnake"},
		{"A", "a"},
		{"AA", "aa"},
		{"AaAa", "aaAa"},
		{"AaaAAa", "aaaaAa"},
		{"HTTPRequest", "httpRequest"},
		{"BatteryLifeValue", "batteryLifeValue"},
	}
	for _, test := range tests {
		have := ToLowerCamelCase(test.input)
		if have != test.want {
			t.Errorf("input=%q:\nhave: %q\nwant: %q", test.input, have, test.want)
		}
	}
}

func TestRandStringLength(t *testing.T) {
	for i := 0; i < 10; i++ {
		lenth := rand.Intn(250)
		v := RandString(lenth, RandStringCharsOnlyNumbers.String())
		assert.Equal(t, lenth, len(v))
	}
}

func TestRandStringRepeatedNumbers(t *testing.T) {
	var results []string
	for i := 0; i < 100; i++ {
		value := ""
		for j := 0; j < 4; j++ {
			v := RandString(2, RandStringCharsOnlyNumbers.String())
			value += v + "-"
			//Necessary for the function n to try to generate in the same nanosecond, thus generating the same number
			time.Sleep(time.Nanosecond * 1)
		}
		results = append(results, value)
		for k := 0; k < len(results)-2; k++ {
			assert.NotEqual(t, results[k], results[i], "Value generated should not be equal to the last generated")
		}
	}
}

// TestRandStringRepeatedString test if the characters are repeteated sequencially
func TestRandStringRepeatedString(t *testing.T) {
	var results []string
	for i := 0; i < 100; i++ {
		value := ""
		for j := 0; j < 4; j++ {
			v := RandString(2, RandStringCharsOnlyLetters.String())
			value += v + "-"
			//Necessary for the function n to try to generate in the same nanosecond, thus generating the same number
			time.Sleep(time.Nanosecond * 1)
		}
		results = append(results, value)
		for k := 0; k < len(results)-2; k++ {
			assert.NotEqual(t, results[k], results[i], "Value generated should not be equal to the last generated")
		}
	}
}

func TestFormatString(t *testing.T) {
	tests := []struct{ input, mask, expected string }{
		{"8499999999", "(##) ####-####", "(84) 9999-9999"},
		{"25740703050", "###.###.###-##", "257.407.030-50"},
		{"4444555566667777", "####.####.####.####", "4444.5555.6666.7777"},
		{"4444555566669879", "****.****.****.####", "****.****.****.9879"},
		{"4444555566669879", "****.####.****.####", "****.5555.****.9879"},
		{"4444555566669879", "####.****.****.####", "4444.****.****.9879"},
		{"59605270", "#####-###", "59605-270"},
		{"06444338000119", "##.###.###/####-##", "06.444.338/0001-19"},
	}
	for _, test := range tests {
		returned := Format(test.input, test.mask)
		if returned != test.expected {
			assert.Equalf(t, test.expected, returned, "input=%q:\nreturned: %q\nwant: %q", test.input, returned, test.expected)
		}
	}
}

func TestGetUriFromUrl(t *testing.T) {
	tests := []struct{ input, expected string }{
		{"http://localhost:7000/api/linksoft/nfe/danfe?chave=24211019843657000149550010000003981197279070", "/api/linksoft/nfe/danfe?chave=24211019843657000149550010000003981197279070"},
		{"https://erp.sigelfe.com/api/linksoft/nfe/danfe?chave=123", "/api/linksoft/nfe/danfe?chave=123"},
		{"http://localhost:7000/teste123", "/teste123"},
		{"teste123", "teste123"},
		{"http://localhost:7000", "/"},
		{"http://localhost:7000/", "/"},
		{"ff", "ff"},
	}
	for _, test := range tests {
		returned := GetUriFromUrl(test.input)
		if returned != test.expected {
			assert.Equalf(t, test.expected, returned, "input=%q:\nreturned: %q\nwant: %q", test.input, returned, test.expected)
		}
	}

}
func TestCutString(t *testing.T) {
	tests := []struct {
		str, inicial, final, esperado string
		n                             int
	}{
		{"TESTE (1234) AAAA", "(", ")", "TESTE  AAAA", 1},
		{"AAAAAAAA BBBBBBB VALOR APROXIMADO DE TRIBUTOS FONTE IBPT", "VALOR APROXIMADO", "FONTE IBPT", "AAAAAAAA BBBBBBB ", 1},
	}
	for _, test := range tests {
		returned := CutString(test.str, test.inicial, test.final, test.n)
		if returned != test.esperado {
			assert.Equal(t, test.esperado, returned)
		}
	}
}

func TestWordCount(t *testing.T) {
	testCases := []struct {
		value    string
		expected int
	}{
		{"Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard", 19},
		{"", 0},
		{"rodrigo rodrigues", 2},
		{"rodrigo", 1},
		{" ", 0},
	}
	for _, v := range testCases {
		got := WordCount(v.value)
		assert.Equal(t, v.expected, got, "", v.value, v.expected, got)
	}
}

func TestOnlyNumbers(t *testing.T) {
	testCases := []struct {
		value    string
		expected string
	}{
		{"123 test", "123"},
		{"AAAAAAAbbb CCC  55587 54848 ddd4a45d4f75r1", "55587548484454751"},
		{"", ""},
		{" ", ""},
		{"1", "1"},
		{"0123456789", "0123456789"},
	}
	for _, v := range testCases {
		got := OnlyNumbers(v.value)
		assert.Equal(t, v.expected, got, "", v.value, v.expected, got)
	}
}

func TestGetDiff(t *testing.T) {
	testCases := []struct {
		origin   string
		final    string
		expected string
	}{
		{
			"o rato roeu a roupa do rei de roma",
			"o rato nao roeu a roupa do rei de roma",
			" (++nao ++)",
		},
		{
			"the mouse gnawed the clothes of the king of rome",
			"the mouse did not gnaw the king of rome's clothes",
			" (~~gnawed the clothes of~~) (++did not gnaw++) (++'s clothes++)",
		},
		{
			"o mar é azul e o ceu está claro",
			"o mar é escuro e o ceu está turvo",
			" (~~azul~~) (++escuro++) (~~clar~~) (++turv++)",
		},
		{
			"the sea is blue and the sky is clear",
			"the sea is dark and the sky is cloudy",
			" (~~blue~~) (++dark++) (~~ear~~) (++oudy++)",
		},
	}
	for _, v := range testCases {
		got := GetDiff(v.origin, v.final)
		assert.Equal(t, v.expected, got, "", v.origin, v.final, got)
	}
}

func TestGetStringBetweenAll(t *testing.T) {
	testCases := []struct {
		input        string
		ini          string
		fin          string
		includeQuery bool
		expected     string
	}{
		{
			"someone (Manuel) is very tired",
			"(",
			")",
			true,
			" (Manuel)",
		},
		{
			"someone (Manuel) is very tired",
			"(",
			")",
			false,
			" Manuel",
		},
		{
			"someone (Manuel) is very tired but the other (Marcos) traveled",
			"(",
			")",
			true,
			" (Manuel) (Marcos)",
		},
		{
			"someone (Manuel) is very tired but the other (Marcos) traveled",
			"(",
			")",
			false,
			" Manuel Marcos",
		},
	}
	for _, v := range testCases {
		got := GetStringBetweenAll(v.input, v.ini, v.fin, v.includeQuery)
		assert.Equal(t, v.expected, got, "", v.expected, got)
	}
}
