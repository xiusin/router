// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pine

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/xiusin/pine/render"
)

type H map[string]interface{}

var engines = map[string]render.IRenderer{}

type Render struct {
	engines map[string]render.IRenderer
	writer  http.ResponseWriter
	tplData H
	charset string

	applied bool
}

const (
	contentTypeJSON = "application/json"
	contentTypeHTML = "text/html"
	contentTypeText = "text/plain"
	contentTypeXML  = "text/xml"
)

func RegisterViewEngine(engine render.IRenderer) {
	engines[engine.Ext()] = engine
}

func newRender(writer http.ResponseWriter, charset string) *Render {

	return &Render{
		engines,
		writer,
		H{},
		charset,
		false,
	}
}

func (c *Render) ContentType(typ string) {
	c.writer.Header().Set("Content-Type", fmt.Sprintf("%s; Charset=%s", typ, c.charset))
}

func (c *Render) reset(writer http.ResponseWriter) {
	c.writer = writer
	c.tplData = H{}
	c.applied = false
}
func (c *Render) JSON(v interface{}) error {
	c.ContentType(contentTypeJSON)

	return responseJson(c.writer, v, "")
}

func (c *Render) Text(v string) error {
	c.ContentType(contentTypeText)

	return c.Bytes([]byte(v))
}

func (c *Render) Bytes(v []byte) error {
	c.ContentType(contentTypeText)

	_, err := c.writer.Write(v)
	return err
}

func (c *Render) HTML(viewPath string) {
	c.ContentType(contentTypeHTML)
	if err := c.engines[filepath.Ext(viewPath)].HTML(c.writer, viewPath, c.tplData); err != nil {
		panic(err)
	}

	c.applied = true
}

func (c *Render) JSONP(callback string, v interface{}) error {
	c.ContentType(contentTypeJSON)

	return responseJson(c.writer, v, callback)
}

func (c *Render) ViewData(key string, val interface{}) {
	c.tplData[key] = val
}

func (c *Render) XML(v interface{}) error {
	c.ContentType(contentTypeXML)

	b, err := xml.MarshalIndent(v, "", " ")
	if err == nil {
		_, err = c.writer.Write(b)
	}

	return err
}

func responseJson(writer io.Writer, v interface{}, callback string) error {
	b, err := json.Marshal(v)
	if err == nil {
		if callback == "" {
			_, err = writer.Write(b)
		} else {
			var ret bytes.Buffer
			ret.Write([]byte(callback))
			ret.Write([]byte("("))
			ret.Write(b)
			ret.Write([]byte(")"))
			_, err = writer.Write(ret.Bytes())
		}
	}
	return err
}
