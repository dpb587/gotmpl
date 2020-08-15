package gotmpl

import "github.com/dpb587/gotmpl/pkg/app"

func NewCommand(app app.Version) *Command {
	return &Command{
		Runtime: NewRuntime(app),
	}
}
