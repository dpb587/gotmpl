package gotmpl

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dpb587/gotmpl/pkg/cliopt"
	"github.com/dpb587/gotmpl/pkg/data"
	gotmplioutil "github.com/dpb587/gotmpl/pkg/ioutil"
	"github.com/dpb587/gotmpl/pkg/template"
	"github.com/pkg/errors"
)

type DataOptions struct {
	JSON cliopt.DataSourceList `long:"json" description:"file or URL to JSON data" value-name:"[KEY=]URI"`
	YAML cliopt.DataSourceList `long:"yaml" description:"file or URL to YAML data" value-name:"[KEY=]URI"`
	Raw  cliopt.DataSourceList `long:"raw" description:"file or URL to raw data" value-name:"[KEY=]URI"`
}

type TemplateOptions struct {
	Template       []string `long:"template" description:"inline template to parse" value-name:"TEMPLATE"`
	NamedTemplates bool     `long:"named-templates" description:"import templates as named blocks based on file name"`
}

type OutputOptions struct {
	Output      cliopt.BlockOutputList `long:"output" description:"write result to file (with optional template block)" value-name:"PATH[=BLOCK]"`
	OutputBlock string                 `long:"output-block" description:"default template block to render" value-name:"BLOCK"`
	OutputType  string                 `long:"output-type" description:"type of output to produce (values: html, text)" default:"text"`
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
	log := c.Runtime.Logger()

	tmpl, err := template.New(c.OutputOptions.OutputType)
	if err != nil {
		return err
	}

	{ // templates
		log.Debugf("loading templates")

		var templateFileCount int

		for _, templateGlob := range c.Args.Templates {
			templateFiles, err := filepath.Glob(templateGlob)
			if err != nil {
				return errors.Wrapf(err, "matching %s", templateGlob)
			}

			if len(templateFiles) == 0 {
				log.Warnf("no templates matched pattern: %s", templateGlob)
			}

			for _, templateFile := range templateFiles {
				buf, err := ioutil.ReadFile(templateFile)
				if err != nil {
					return errors.Wrapf(err, "reading %s", templateFile)
				}

				if c.TemplateOptions.NamedTemplates {
					tmpl, err = tmpl.ParseNamed(filepath.Base(templateFile), string(buf))
				} else {
					tmpl, err = tmpl.Parse(string(buf))
				}

				if err != nil {
					return errors.Wrapf(err, "parsing %s", templateFile)
				}

				log.Debugf("loaded template (%s)", templateFile)

				templateFileCount++
			}
		}

		var templateInlineCount int

		if inline := c.TemplateOptions.Template; len(inline) > 0 {
			for templateRawIdx, templateRaw := range inline {
				tmpl, err = tmpl.Parse(templateRaw)
				if err != nil {
					return errors.Wrapf(err, "parsing inline template %d", templateRawIdx)
				}

				templateInlineCount++
			}

		}

		log.Infof("loaded templates (files: %d, inlines: %d)", templateFileCount, templateInlineCount)

		if templateFileCount+templateInlineCount == 0 {
			return errors.New("no templates found")
		}
	}

	var tmplData interface{}

	{ // data
		var sources []*data.Source
		httpClient := c.Runtime.NewHTTPClient()

		for parser, rawSources := range map[data.Parser]cliopt.DataSourceList{
			data.JSON: c.DataOptions.JSON,
			data.YAML: c.DataOptions.YAML,
			data.Raw:  c.DataOptions.Raw,
		} {
			for _, rawSource := range rawSources {
				var source data.Opener

				if strings.HasPrefix(rawSource.URI, "http://") || strings.HasPrefix(rawSource.URI, "https://") {
					source = data.NewHTTPOpener(httpClient, rawSource.URI)
				} else if rawSource.URI == "/dev/stdin" || rawSource.URI == "-" {
					source = data.NewReaderOpener(os.Stdin)
				} else {
					source = data.NewFileOpener(rawSource.URI)
				}

				sources = append(sources, data.NewSource(rawSource.Key, source, parser))
			}
		}

		if len(sources) > 0 {
			log.Debugf("loading data sources")

			dataWrap := data.NewData(sources...)

			tmplData, err = dataWrap.Data()
			if err != nil {
				return errors.Wrap(err, "loading data")
			}

			log.Info("loaded data sources")
		}
	}

	{ // output
		outps := c.OutputOptions.Output

		if len(outps) == 0 {
			outps = append(outps, cliopt.BlockOutput{Block: "", Path: "/dev/stdout"})
		}

		for outpIdx, outp := range outps {
			if outp.Block == "" {
				outp.Block = c.OutputOptions.OutputBlock
			}

			log.Debugf("rendering %s from block %s", outp.PathLabel(), outp.BlockLabel(tmpl.Name()))

			out, closer, err := outp.Open()
			if err != nil {
				return errors.Wrapf(err, "rendering output %d", outpIdx)
			}

			outw := gotmplioutil.NewWriterCounter(out)

			if v := outp.Block; v != "" {
				err = tmpl.ExecuteTemplate(outw, v, tmplData)
			} else {
				err = tmpl.Execute(outw, tmplData)
			}

			if closer != nil {
				closer()
			}

			if err != nil {
				return errors.Wrap(err, "rendering")
			}

			log.Infof("rendered %s from block %s (bytes: %d)", outp.PathLabel(), outp.BlockLabel(tmpl.Name()), outw.WriteLength())
		}
	}

	return nil
}
