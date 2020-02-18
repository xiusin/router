package main

import (
	"fmt"
	"github.com/xiusin/pine"
)

func main() {
	app := pine.New()

	app.GET("/panic", func(ctx *pine.Context) {
		panic("服务错误")
	})

	app.SetRecoverHandler(func(ctx *pine.Context) {
		if err := recover(); err != nil {
			ctx.Render().Text(fmt.Sprintf("recover函数必须放到recover的判断里: %s", err))
		}
	})

	// 使用debug组件替换默认recover函数
	//app.SetRecoverHandler(debug.Recover(app))

	app.Run(pine.Addr(":9528"))
}
