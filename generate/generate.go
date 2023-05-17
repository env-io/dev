package generate

import (
	"fmt"
	"os"
	"strings"

	"github.com/env-io/dev/core"
)

func DirExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func ContainsString(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}

func AskForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		core.LogError(err.Error())
	}
	ok := []string{"y", "Y", "yes", "Yes", "YES"}
	notOk := []string{"n", "N", "no", "No", "NO"}
	if ContainsString(ok, response) {
		return true
	} else if ContainsString(notOk, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return AskForConfirmation()
	}
}

func FileReader(file string) (f *os.File, err error) {
	if DirExist(file) {
		core.LogInfo(fmt.Sprintf("%v is exist, \n Do you want to overwrite it ? Yes or No?", file))
		if AskForConfirmation() {
			if f, err = os.OpenFile(file, os.O_RDWR|os.O_TRUNC, 0666); err != nil {
				core.LogError(err.Error())
				return
			}
		} else {
			core.LogInfo("Skip creating file.")
			return
		}
	} else {
		if f, err = os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0666); err != nil {
			core.LogError(err.Error())
			return
		}
	}

	return
}

func StubReplaces(content string, tpl *core.StubTemplate) string {
	content = strings.Replace(content, "{{ProjectPath}}", tpl.ProjectPath, -1)
	content = strings.Replace(content, "{{PackagePath}}", tpl.PackagePath, -1)
	content = strings.Replace(content, "{{PackageName}}", core.ToLowerCamelCase(tpl.PackageName), -1)
	content = strings.Replace(content, "{{ModuleName}}", tpl.ModuleName, -1)
	content = strings.Replace(content, "{{ModelName}}", tpl.ModelName, -1)
	content = strings.Replace(content, "{{ModelNameSingular}}", tpl.ModelNameSingular, -1)
	content = strings.Replace(content, "{{ModelNamePlural}}", tpl.ModelNamePlural, -1)
	content = strings.Replace(content, "{{TableName}}", tpl.TableName, -1)

	return content
}

func WriteFile(file *os.File, content string, tpl *core.StubTemplate) {
	if tpl != nil {
		content = StubReplaces(content, tpl)
	}

	if _, err := file.WriteString(content); err != nil {
		core.LogError(fmt.Sprintf("Could not write file %s\n%s", file.Name(), err.Error()))
		os.Exit(2)
	}

	file.Close()
}
