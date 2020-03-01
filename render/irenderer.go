// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"io"
)

type IRenderer interface {
	AddFunc(string, interface{})
	HTML(writer io.Writer, name string, binding map[string]interface{}) error
	Ext() string
}
