package template

import (
	"io"
	"text/template"
)

type textTemplate struct {
	tmpl *template.Template
}

var _ Template = &textTemplate{}

func (t *textTemplate) Parse(text string) (Template, error) {
	tmpl, err := t.tmpl.Parse(text)
	if err != nil {
		return nil, err
	}

	res := &textTemplate{
		tmpl: tmpl,
	}

	return res, nil
}

func (t *textTemplate) ParseFiles(filenames ...string) (Template, error) {
	tmpl, err := t.tmpl.ParseFiles(filenames...)
	if err != nil {
		return nil, err
	}

	res := &textTemplate{
		tmpl: tmpl,
	}

	return res, nil
}

func (t *textTemplate) Funcs(funcMap FuncMap) Template {
	return &textTemplate{
		tmpl: t.tmpl.Funcs(template.FuncMap(funcMap)),
	}
}

func (t *textTemplate) Execute(wr io.Writer, data interface{}) error {
	return t.tmpl.Execute(wr, data)
}

func (t *textTemplate) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return t.tmpl.ExecuteTemplate(wr, name, data)
}
