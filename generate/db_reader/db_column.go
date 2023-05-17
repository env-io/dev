package dbReader

import (
	"fmt"
	"strings"

	"github.com/env-io/dev/core"
)

type Column struct {
	Name string
	Type string
	Tag  *OrmTag
}

func (col *Column) String() string {
	if strings.HasSuffix(col.Name, "Id") {
		col.Name = core.RightTrim(col.Name, "Id") + "ID"
	}
	return fmt.Sprintf("%s %s %s", col.Name, col.Type, col.Tag.String())
}

type ForeignKey struct {
	Name      string
	RefSchema string
	RefTable  string
	RefColumn string
	Column    *Column
}
