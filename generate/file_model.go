package generate

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/butbetter-id/dev/core"
	dbReader "github.com/butbetter-id/dev/generate/db_reader"
	"github.com/butbetter-id/dev/generate/stubs"
)

func FileModel(driver string, conn string, selectedTables string, tpl *core.StubTemplate) {

	var tables map[string]bool
	if selectedTables != "" {
		tables = make(map[string]bool)
		for _, v := range strings.Split(selectedTables, ",") {
			tables[v] = true
		}
	}

	db, err := sql.Open(driver, conn)
	if err != nil {
		core.LogError("Could not connect to database ")
		core.LogError(fmt.Sprintf("using: %s, %s, %s", driver, conn, err.Error()))
		os.Exit(2)
	}
	defer db.Close()

	if trans, ok := dbReader.DBDriver[driver]; ok {
		core.LogInfo("")
		core.LogInfo("Generating model file ...")
		core.LogInfo("--------------------------------------")

		tableNames := trans.GetTableNames(db)
		dbTables := dbReader.GetTableObjects(tableNames, db, trans)

		makeModels(dbTables, tables, tpl.AppPath, tpl)
	} else {
		core.LogError(fmt.Sprintf("%s database is not supported yet.", driver))
		os.Exit(2)
	}
}

func makeModels(tables []*dbReader.Table, selectedTables map[string]bool, modelPath string, tpl *core.StubTemplate) {
	for _, tb := range tables {
		if selectedTables != nil {
			if _, selected := selectedTables[tb.Name]; !selected {
				continue
			}
		}

		filename := dbReader.GetFileName(tb.Name)
		file := path.Join(modelPath, filename+".go")
		f, err := FileReader(file)
		if err != nil {
			continue
		}

		template := stubs.Model
		template = strings.Replace(template, "{{modelStruct}}", tb.String(), 1)

		var tPkg string
		if tb.ImportTimePkg {
			tPkg = "\"time\"\n"
		}

		template = strings.Replace(template, "{{timePkg}}", tPkg, -1)

		tpl.ModelName = core.ToCamelCase(tb.Name)
		tpl.ModelNameSingular = strings.TrimSuffix(tpl.ModelName, "s")
		tpl.TableName = tb.Name
		WriteFile(f, template, tpl)
		core.FormatSourceCode(f.Name())
		core.LogInfo(fmt.Sprintf("%-20s => \t\t%s", "model", file))
	}
}
