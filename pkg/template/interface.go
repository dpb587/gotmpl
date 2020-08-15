package template

import (
	"io"
)

type FuncMap map[string]interface{}

type Template interface {
	Parse(string) (Template, error)
	ParseFiles(...string) (Template, error)

	Funcs(funcMap FuncMap) Template

	Execute(io.Writer, interface{}) error
	ExecuteTemplate(io.Writer, string, interface{}) error
}
