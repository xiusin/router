// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pine

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"path/filepath"

	"github.com/xiusin/pine/render"
)

type H map[string]interface{}

var engines = map[string]render.AbstractRenderer{}

type Render struct {
	engines map[string]render.AbstractRenderer
	writer  http.ResponseWriter
	tplData H

	applied bool
}

const (
	ContentTypeJSON = "application/json; charset=utf-8"
	ContentTypeHTML = "text/html; charset=utf-8"
	ContentTypeText = "text/plain; charset=utf-8"
	ContentTypeXML  = "text/xml; charset=utf-8"
)

func RegisterViewEngine(engine render.AbstractRenderer) {
	if engine == nil {
		panic("engine can not be nil")
	}
	engines[engine.Ext()] = engine
}

func newRender(writer http.ResponseWriter) *Render {
	return &Render{
		engines,
		writer,
		nil,
		false,
	}
}

func (c *Render) ContentType(typ string) {
	c.writer.Header().Set("Content-Type", typ)
}

func (c *Render) reset(writer http.ResponseWriter) {
	c.writer = writer
	for k := range c.tplData {
		delete(c.tplData, k)
	}
	c.applied = false
}
func (c *Render) JSON(v interface{}) error {
	c.writer.Header().Set("Content-Type", ContentTypeJSON)

	return responseJson(c.writer, v, "")
}

func (c *Render) Text(v string) error {
	return c.Bytes([]byte(v))
}

func (c *Render) Bytes(v []byte) error {
	_, err := c.writer.Write(v)
	return err
}

func (c *Render) HTML(viewPath string) {
	c.writer.Header().Set("Content-Type", ContentTypeHTML)

	if err := c.engines[filepath.Ext(viewPath)].HTML(c.writer, viewPath, c.tplData); err != nil {
		panic(err)
	}

	c.applied = true
}

func (c *Render) GetEngine(ext string) render.AbstractRenderer {
	return c.engines[ext]
}

func (c *Render) JSONP(callback string, v interface{}) error {
	c.writer.Header().Set("Content-Type", ContentTypeJSON)

	return responseJson(c.writer, v, callback)
}

func (c *Render) ViewData(key string, val interface{}) {
	if c.tplData == nil {
		c.tplData = H{}
	}
	c.tplData[key] = val
}

func (c *Render) GetViewData() map[string]interface{} {
	return c.tplData
}

func (c *Render) XML(v interface{}) error {
	c.writer.Header().Set("Content-Type", ContentTypeXML)

	b, err := xml.MarshalIndent(v, "", " ")
	if err == nil {
		_, err = c.writer.Write(b)
	}

	return err
}

func responseJson(writer io.Writer, v interface{}, callback string) error {
	b, err := json.Marshal(v)
	if err == nil {
		if len(callback) == 0 {
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
