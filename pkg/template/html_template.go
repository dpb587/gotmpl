package template

import (
	"html/template"
	"io"
)

type htmlTemplate struct {
	tmpl *template.Template
}

var _ Template = &htmlTemplate{}

func (t *htmlTemplate) Parse(text string) (Template, error) {
	tmpl, err := t.tmpl.Parse(text)
	if err != nil {
		return nil, err
	}

	res := &htmlTemplate{
		tmpl: tmpl,
	}

	return res, nil
}

func (t *htmlTemplate) ParseFiles(filenames ...string) (Template, error) {
	tmpl, err := t.tmpl.ParseFiles(filenames...)
	if err != nil {
		return nil, err
	}

	res := &htmlTemplate{
		tmpl: tmpl,
	}

	return res, nil
}

func (t *htmlTemplate) Funcs(funcMap FuncMap) Template {
	return &htmlTemplate{
		tmpl: t.tmpl.Funcs(template.FuncMap(funcMap)),
	}
}

func (t *htmlTemplate) Execute(wr io.Writer, data interface{}) error {
	return t.tmpl.Execute(wr, data)
}

func (t *htmlTemplate) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return t.tmpl.ExecuteTemplate(wr, name, data)
}
