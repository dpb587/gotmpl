package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/dpb587/gotmpl/cmd/gotmpl"
	"github.com/dpb587/gotmpl/pkg/app"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

var (
	appName                        = "gotmpl"
	appOrigin                      = "github.com/dpb587/gotmpl"
	appSemver, appCommit, appBuilt string
)

func main() {
	v := app.MustVersion(appName, appSemver, appCommit, appBuilt)
	cmd := gotmpl.NewCommand(v)

	parser := flags.NewParser(cmd, flags.PassDoubleDash)

	fatal := func(err error) {
		if debug, _ := strconv.ParseBool(os.Getenv("DEBUG")); debug {
			panic(err)
		}

		fmt.Fprintf(os.Stderr, "%s: error: %s\n", parser.Command.Name, err)

		os.Exit(1)
	}

	_, err := parser.Parse()
	if err != nil {
		fatal(err)
	} else if cmd.Runtime.Help {
		helpBuf := &bytes.Buffer{}
		parser.WriteHelp(helpBuf)
		help := helpBuf.String()

		// join conventional paren groups
		help = strings.Replace(help, ") (", "; ", -1)

		fmt.Print(help)
		fmt.Printf("\n")

		return
	} else if cmd.Runtime.Version != nil {
		ver, err := semver.NewVersion(v.Semver)
		if err != nil {
			fatal(errors.Wrap(err, "parsing application version"))
		}

		if !cmd.Runtime.Quiet {
			app.WriteVersion(os.Stdout, os.Args[0], v, len(cmd.Runtime.Verbose))
		}

		if !cmd.Runtime.Version.Constraint.Check(ver) {
			if cmd.Runtime.Quiet {
				os.Exit(1)
			}

			fatal(errors.Wrapf(fmt.Errorf("constraint not met: %s", cmd.Runtime.Version.Constraint.RawValue), "verifying application version"))
		}

		return
	}

	if err = cmd.Execute(nil); err != nil {
		fatal(err)
	}
}
