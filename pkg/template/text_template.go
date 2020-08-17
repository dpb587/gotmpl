package template

import (
	"io"
	"text/template"

	"github.com/pkg/errors"
)

type textTemplate struct {
	tmpl *template.Template
}

var _ Template = &textTemplate{}

func (t *textTemplate) Name() string {
	return t.tmpl.Name()
}

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

func (t *textTemplate) ParseNamed(name string, text string) (Template, error) {
	tmpl, err := template.New(name).Parse(text)
	if err != nil {
		return nil, err
	}

	tmpl, err = t.tmpl.AddParseTree(name, tmpl.Tree)
	if err != nil {
		return nil, errors.Wrapf(err, "adopting %s", name)
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
