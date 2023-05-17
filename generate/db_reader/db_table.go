package dbReader

import (
	"fmt"
	"strings"

	"github.com/env-io/dev/core"
)

type Table struct {
	Name          string
	Pk            string
	Uk            []string
	Fk            map[string]*ForeignKey
	Columns       []*Column
	ImportTimePkg bool
}

func (tb *Table) String() string {
	rv := fmt.Sprintf("type %s struct {\n", core.ToCamelCase(tb.Name))
	for _, v := range tb.Columns {
		rv += v.String() + "\n"
	}
	rv += "}\n"
	return rv
}

func (tb *Table) MarshalColumn() string {
	var colMarshal []string

	colMarshal = append(colMarshal, fmt.Sprintf("%s %s %s", "ID", "string", "`json:\"id\"`"))
	for col, fk := range tb.Fk {
		cname := core.ToCamelCase(fk.Name)
		if strings.HasSuffix(cname, "Id") {
			cname = core.RightTrim(cname, "Id") + "ID"
		} else {
			cname = cname + "ID"
			col = col + "_id"
		}
		colMarshal = append(colMarshal, fmt.Sprintf("%s %s %s", cname, "string", fmt.Sprintf("`json:\"%s\"`", col)))
	}
	return strings.Join(colMarshal, "\n")
}
