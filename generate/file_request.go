package generate

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/env-io/dev/core"
	"github.com/env-io/dev/generate/stubs"
)

func FileRequest(name string, tpl *core.StubTemplate) {
	fileHandler := path.Join(tpl.AppPath, fmt.Sprintf("request_%s.go", name))
	f, err := FileReader(fileHandler)
	if err != nil {
		os.Exit(2)
	}

	rtemplate := stubs.RequestHeader
	rtemplate += strings.Replace(stubs.RequestStruct, "{{RequestName}}", core.ToLower(name), -1)
	WriteFile(f, rtemplate, tpl)
	core.FormatSourceCode(f.Name())
	core.LogInfo(fmt.Sprintf("%-20s => \t\t%s", "request", f.Name()))
}
