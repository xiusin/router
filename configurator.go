// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pine

import "github.com/xiusin/pine/sessions/cookie_transcoder"

type Configuration struct {
	maxMultipartMemory        int64
	serverName                string
	withoutStartupLog         bool
	autoParseControllerResult bool
	autoParseForm             bool
	useCookie                 bool
	CookieTranscoder          cookie_transcoder.AbstractCookieTranscoder
}

type AbstractReadonlyConfiguration interface {
	GetServerName() string
	GetAutoParseForm() bool
	GetUseCookie() bool
	GetMaxMultipartMemory() int64
	GetAutoParseControllerResult() bool
	GetCookieTranscoder() cookie_transcoder.AbstractCookieTranscoder
}

type Configurator func(o *Configuration)

func WithServerName(srvName string) Configurator {
	return func(o *Configuration) {
		o.serverName = srvName
	}
}

func WithCookieTranscoder(transcoder cookie_transcoder.AbstractCookieTranscoder) Configurator {
	return func(o *Configuration) {
		o.CookieTranscoder = transcoder
	}
}

func WithCookie(open bool) Configurator {
	return func(o *Configuration) {
		o.useCookie = open
	}
}

func WithAutoParseForm(autoParseForm bool) Configurator {
	return func(o *Configuration) {
		o.autoParseForm = autoParseForm
	}
}

func WithMaxMultipartMemory(mem int64) Configurator {
	return func(o *Configuration) {
		o.maxMultipartMemory = mem
	}
}

func WithoutStartupLog(hide bool) Configurator {
	return func(o *Configuration) {
		o.withoutStartupLog = hide
	}
}

func WithAutoParseControllerResult(auto bool) Configurator {
	return func(o *Configuration) {
		o.autoParseControllerResult = auto
	}
}

func (c *Configuration) GetServerName() string {
	return c.serverName
}

func (c *Configuration) GetUseCookie() bool {
	return c.useCookie
}

func (c *Configuration) GetAutoParseForm() bool {
	return c.autoParseForm
}

func (c *Configuration) GetAutoParseControllerResult() bool {
	return c.autoParseControllerResult
}

func (c *Configuration) GetCookieTranscoder() cookie_transcoder.AbstractCookieTranscoder {
	return c.CookieTranscoder
}

func (c *Configuration) GetMaxMultipartMemory() int64 {
	return c.maxMultipartMemory
}
