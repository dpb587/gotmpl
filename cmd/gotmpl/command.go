package gotmpl

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dpb587/gotmpl/pkg/cliopt"
	"github.com/dpb587/gotmpl/pkg/data"
	"github.com/dpb587/gotmpl/pkg/template"
	"github.com/pkg/errors"
)

type DataOptions struct {
	JSONFiles cliopt.KeyValueList `long:"json-file" description:"path to a JSON file for data" value-name:"[KEY=]PATH"`
	JSONURLs  cliopt.KeyValueList `long:"json-url" description:"download URL to a JSON file for data" value-name:"[KEY=]URL"`

	YAMLFiles cliopt.KeyValueList `long:"yaml-file" description:"path to a YAML file for data" value-name:"[KEY=]PATH"`
	YAMLURLs  cliopt.KeyValueList `long:"yaml-url" description:"download URL to a YAML file for data" value-name:"[KEY=]URL"`

	TextFiles cliopt.KeyValueList `long:"text-file" description:"path to a file for raw data" value-name:"[KEY=]PATH"`
	TextURLs  cliopt.KeyValueList `long:"text-url" description:"download URL to a file for raw data" value-name:"[KEY=]URL"`
}

type TemplateOptions struct {
	Template []string `long:"template" description:"inline template to parse" value-name:"TEMPLATE"`
}

type OutputOptions struct {
	OutputType string              `long:"output-type" description:"type of output to produce (values: html, text)" default:"text"`
	Output     cliopt.KeyValueList `long:"output" description:"write result to file" value-name:"[BLOCK=]PATH"`
}

type Command struct {
	*Runtime         `group:"Runtime Options"`
	*TemplateOptions `group:"Template Options"`
	*DataOptions     `group:"Data Options"`
	*OutputOptions   `group:"Output Options"`
	Args             CommandArgs `positional-args:"true"`
}

type CommandArgs struct {
	Templates []string `positional-arg-name:"TEMPLATE-GLOB" description:"template file(s)"`
}

func (c *Command) Execute(_ []string) error {
	tmpl, err := template.New(c.OutputOptions.OutputType)
	if err != nil {
		return err
	}

	{ // templates
		var templateCount int

		for _, templateGlob := range c.Args.Templates {
			templateFiles, err := filepath.Glob(templateGlob)
			if err != nil {
				return errors.Wrapf(err, "matching %s", templateGlob)
			}

			for _, templateFile := range templateFiles {
				buf, err := ioutil.ReadFile(templateFile)
				if err != nil {
					return errors.Wrapf(err, "reading %s", templateFile)
				}

				tmpl, err = tmpl.Parse(string(buf))
				if err != nil {
					return errors.Wrapf(err, "parsing %s", templateFile)
				}

				templateCount++
			}
		}

		for templateRawIdx, templateRaw := range c.TemplateOptions.Template {
			tmpl, err = tmpl.Parse(templateRaw)
			if err != nil {
				return errors.Wrapf(err, "parsing inline template %d", templateRawIdx)
			}

			templateCount++
		}

		if templateCount == 0 {
			return errors.New("missing template")
		}
	}

	var tmplData interface{}

	{ // data
		var sources []*data.Source
		httpClient := c.Runtime.NewHTTPClient()

		for _, opt := range c.DataOptions.JSONFiles {
			sources = append(
				sources,
				data.NewSource(
					opt.Key,
					data.NewFileOpener(opt.Value),
					data.JSON,
				),
			)
		}

		for _, opt := range c.DataOptions.JSONURLs {
			sources = append(
				sources,
				data.NewSource(
					opt.Key,
					data.NewURLOpener(httpClient, opt.Value),
					data.JSON,
				),
			)
		}

		for _, opt := range c.DataOptions.YAMLFiles {
			sources = append(
				sources,
				data.NewSource(
					opt.Key,
					data.NewFileOpener(opt.Value),
					data.YAML,
				),
			)
		}

		for _, opt := range c.DataOptions.YAMLURLs {
			sources = append(
				sources,
				data.NewSource(
					opt.Key,
					data.NewURLOpener(httpClient, opt.Value),
					data.YAML,
				),
			)
		}

		for _, opt := range c.DataOptions.TextFiles {
			sources = append(
				sources,
				data.NewSource(
					opt.Key,
					data.NewFileOpener(opt.Value),
					data.Raw,
				),
			)
		}

		for _, opt := range c.DataOptions.TextURLs {
			sources = append(
				sources,
				data.NewSource(
					opt.Key,
					data.NewURLOpener(httpClient, opt.Value),
					data.Raw,
				),
			)
		}

		dataWrap := data.NewData(sources...)

		tmplData, err = dataWrap.Data()
		if err != nil {
			return errors.Wrap(err, "loading data")
		}
	}

	{ // output
		outps := c.OutputOptions.Output

		if len(outps) == 0 {
			outps = append(outps, cliopt.KeyValue{Key: "", Value: "-"})
		}

		for _, outp := range outps {
			outBlock := outp.Key
			outPath := outp.Value

			if outp.Key != "" {
				outBlock = outp.Value
				outPath = outp.Key
			}

			var out io.Writer

			if outPath == "-" {
				out = os.Stdout
			} else {
				fh, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
				if err != nil {
					return errors.Wrapf(err, "opening %s", outPath)
				}

				defer fh.Close() // TODO immediate

				out = fh
			}

			if outBlock == "" {
				err = tmpl.Execute(out, tmplData)
			} else {
				err = tmpl.ExecuteTemplate(out, outBlock, tmplData)
			}

			if err != nil {
				return errors.Wrap(err, "rendering")
			}
		}
	}

	return nil
}
