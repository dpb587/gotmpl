package template

import (
	"fmt"
	html "html/template"
	text "text/template"

	"github.com/Masterminds/sprig"
)

const templateName = "gotmpl"

func New(mode string) (Template, error) {
	var res Template

	if mode == "html" {
		res = &htmlTemplate{
			tmpl: html.New(templateName),
		}
	} else if mode == "text" {
		res = &textTemplate{
			tmpl: text.New(templateName),
		}
	} else {
		return nil, fmt.Errorf("unsupported mode: %s", mode)
	}

	res = res.Funcs(FuncMap(sprig.FuncMap()))

	return res, nil
}
