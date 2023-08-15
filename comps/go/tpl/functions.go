package tpl

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/waykiss/wkcomps/file"
	"github.com/waykiss/wkcomps/html"
	"github.com/waykiss/wkcomps/number"
	"github.com/waykiss/wkcomps/str"
	"github.com/waykiss/wkcomps/validation"
	"github.com/waykiss/wkcomps/xml"
)

var renderTemplateFuncMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
	"formatDinheiro": func(v float64, p int) string {
		return number.GetFormattedMoneyBRL(v, p)
	},
	"isLast": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	},
	"sum": func(value1, value2 float64) float64 {
		return value1 + value2
	},
	"sub": func(value1, value2 float64) float64 {
		return value1 - value2
	},
	"multi": func(value1, value2 float64) float64 {
		return value1 * value2
	},
	"formatDinheiroSemPrefixo": func(v float64, p int) string {
		return number.GetFormattedMoneyBRLWithNoPrefix(v, p)
	},
	"formatNumero": func(v float64, emptyZero bool) string {
		if v == 0 && emptyZero {
			return ""
		}
		return number.GetFormattedFloat(v)
	},
	"formatTelefone": func(v string) string {
		return str.Format(v, "(##) #####-####")
	},
	"formatCep": func(v string) string {
		return str.Format(v, "#####-###")
	},
	"formatCpfCnpj": func(v string) string {
		if validation.IsCpfValid(v) {
			return str.Format(v, str.FormatCpf)
		}
		return str.Format(v, str.FormatCnpj)
	},
	"formatDataHoraPadrao": func(v time.Time) string {
		if v.IsZero() {
			return ""
		}
		return v.Format("02/01/2006 15:04:05")
	},
	"formatData": func(v time.Time) string {
		if v.IsZero() {
			return ""
		}
		return v.Format("02/01/2006")
	},
	"formatDataHora": func(v time.Time, format string) string {
		if v.IsZero() {
			return ""
		}
		return v.Format(format)
	},
}

// GetDefaultTemplateFunctions retorna as funcoes padrao ja pra usar nos templates
func GetDefaultTemplateFunctions() template.FuncMap {
	return renderTemplateFuncMap
}

// RenderTemplate renderiza um template baseado em uma lista de arquivos passado por parametro, essa renderizacao tem suporte
// a nested template, ou seja, um template que precisa de outro e a funcao trata isso também
func RenderTemplate(templates embed.FS, templateName string, data interface{}, functions *template.FuncMap) (r string, err error) {
	dir := filepath.Dir(templateName)
	dir = strings.ReplaceAll(dir, "\\", "/")
	tpl, err := template.New("templates").Funcs(renderTemplateFuncMap).ParseFS(templates, fmt.Sprintf("%s/*", dir))
	if err != nil {
		return
	}
	if functions != nil {
		_ = tpl.Funcs(*functions)
	}
	tpl, err = tpl.ParseFS(templates, templateName)
	if err != nil {
		return
	}

	buf := new(bytes.Buffer)
	err = tpl.ExecuteTemplate(buf, file.GetFileNameWithExtensionFromPath(templateName), data)
	r = minifyTemplate(buf.String(), templateName)
	return
}

// RenderTemplateString renderiza um template a partir de uma string como parametro(templateContent)
// se for XML ou XML já entrega minificado
func RenderTemplateString(templateName, templateContent string, data interface{}, functions *template.FuncMap) (r string, err error) {
	tpl := template.New("templates").Funcs(renderTemplateFuncMap)
	if functions != nil {
		_ = tpl.Funcs(*functions)
	}
	tpl, err = tpl.Parse(templateContent)
	if err != nil {
		return
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		return
	}
	r = minifyTemplate(buf.String(), templateName)
	return
}

// minifyTemplate minifica um template baseado na extensao do nome do arquivo
func minifyTemplate(templateNameString string, fileName string) (r string) {
	switch file.GetFileExtensionFromPath(fileName) {
	case ".html", ".gohtml":
		r = html.MinifyHtml(templateNameString)
	case ".xml":
		if strings.Contains(templateNameString, "&") && !strings.Contains(templateNameString, "&amp;") {
			templateNameString = strings.ReplaceAll(templateNameString, "&", "&amp;")
		}
		r = xml.MinifyXml(templateNameString)
		r = file.RemoveXmlVersion(r)
	default:
		r = templateNameString
	}
	return
}
