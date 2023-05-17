// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stubs

var RequestHeader = `
package {{PackageName}}

import (
	"github.com/env-io/validate"
)
`

var RequestStruct = `
type {{RequestName}}Request struct {}

func (r *{{RequestName}}Request) Validate() *validate.Response {
	o := validate.NewResponse()

	return o
}

func (r *{{RequestName}}Request) Messages() map[string]string {
	return map[string]string{}
}
`
