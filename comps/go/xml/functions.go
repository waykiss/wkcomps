package xml

import (
	"bytes"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
	"regexp"
	"strings"
)

func MinifyXml(v string) string {
	m := minify.New()
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	r := strings.NewReader(v)
	buf := &bytes.Buffer{}
	err := m.Minify("text/xml", buf, r)
	if err != nil {
		return "Erro ao minificar XML"
	}
	return buf.String()
}
