// Copyright 2019 Bloomberg Finance L.P.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package docket

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/fatih/color"
)

func writeHelp(out io.Writer) {
	bold := func(s string) string {
		return color.New(color.Bold).Sprint(s)
	}

	tmpl := template.New("help").Funcs(template.FuncMap{"var": bold})
	err := template.Must(tmpl.Parse(`
Help for using docket:

  {{ var "DOCKET_MODE" }}
    To use docket, set this to the name of the mode to use.

Optional environment variables:

  {{ var "DOCKET_DOWN" }} (default off)
    If non-empty, docket will run 'docker-compose down' after each suite.

  {{ var "DOCKET_PULL" }} (default off)
    If non-empty, docket will run 'docker-compose pull' before each suite.

`[1:])).Execute(out, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to Execute help template: %v", err))
	}
}

//nolint:gochecknoinits // This seems to be the only way to get a flag added to go test's help.
func init() {
	// We register a flag to get it shown in the default usage.
	//
	// We don't actually use the parsed flag value, though, since that would require us to call
	// flag.Parse() here. If we call flag.Parse(), then higher-level libraries can't easily add
	// their own flags, since testing's t.Run() will not re-run flag.Parse() if the flags have
	// already been parsed.
	//
	// Instead, we simply look for our flag text in os.Args.

	flag.Bool("help-docket", false, "get help on docket")

	for _, arg := range os.Args {
		if arg == "-help-docket" || arg == "--help-docket" {
			writeHelp(os.Stderr)

			const helpExitCode = 2 // This matches what 'go test -h' returns.
			os.Exit(helpExitCode)
		}
	}
}
