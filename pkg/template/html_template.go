package template

import (
	"html/template"
	"io"

	"github.com/pkg/errors"
)

type htmlTemplate struct {
	tmpl *template.Template
}

var _ Template = &htmlTemplate{}

func (t *htmlTemplate) Name() string {
	return t.tmpl.Name()
}

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

func (t *htmlTemplate) ParseNamed(name string, text string) (Template, error) {
	tmpl, err := template.New(name).Parse(text)
	if err != nil {
		return nil, err
	}

	tmpl, err = t.tmpl.AddParseTree(name, tmpl.Tree)
	if err != nil {
		return nil, errors.Wrapf(err, "adopting %s", name)
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
