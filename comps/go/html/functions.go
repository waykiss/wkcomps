package html

import (
	"bytes"
	"github.com/waykiss/wkcomps/str"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/svg"
	"regexp"
	"strings"
)

func MinifyHtml(v string) string {
	m := minify.New()
	m.AddFunc("image/svg+html", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]html$"), html.Minify)
	r := strings.NewReader(v)
	buf := &bytes.Buffer{}
	err := m.Minify("text/html", buf, r)
	if err != nil {
		return "Erro ao minificar HTML"
	}
	return str.RemoveExtraSpaces(buf.String())
}
