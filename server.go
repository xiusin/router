// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pine

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"strings"
)

type ServerHandler func(*Application) error

const zeroIP = "0.0.0.0"
const defaultAddressWithPort = zeroIP + ":9528"

func (a *Application) printSetupInfo(addr string) {
	if strings.HasPrefix(addr, ":") {
		addr = fmt.Sprintf("%s%s", a.hostname, addr)
	}
	color.Green.Println(logo)
	color.Red.Println("\nServer now listening on: " + addr)
}

func Addr(addr string) ServerHandler {
	return func(a *Application) error {
		s := &fasthttp.Server{}
		s.Logger = Logger()
		s.Name = a.configuration.serverName
		if len(addr) == 0 {
			addr = defaultAddressWithPort
		}
		addrInfo := strings.Split(addr, ":")
		if len(addrInfo[0]) == 0 {
			addrInfo[0] = zeroIP
		}
		a.hostname = addrInfo[0]
		if !a.configuration.withoutStartupLog {
			a.printSetupInfo(addr)
		}
		a.quitCh = make(chan os.Signal)
		if a.configuration.maxMultipartMemory > 0 {
			s.MaxRequestBodySize = int(a.configuration.maxMultipartMemory)
		}
		s.Handler = func(ctx *fasthttp.RequestCtx) {
			c := a.pool.Acquire().(*Context)
			c.beginRequest(ctx)
			defer a.pool.Release(c)
			defer func() { c.RequestCtx = nil }()
			defer c.endRequest(a.recoverHandler)
			a.handle(c)
		}
		signal.Notify(a.quitCh, os.Interrupt, os.Kill)
		go a.gracefulShutdown(s, a.quitCh)
		return s.ListenAndServe(addr)
	}
}
