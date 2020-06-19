package routes

import (
	"genids/controls/api"
	"github.com/kataras/iris/v12"
)

// 定义500错误处理函数
func err500(ctx iris.Context) {
	_, _ = ctx.WriteString("CUSTOM 500 ERROR")
}

// 定义404错误处理函数
func err404(ctx iris.Context) {
	_, _ = ctx.WriteString("CUSTOM 404 ERROR")
}

// 注入路由
func InnerRoute(app *iris.Application) {
	app.OnErrorCode(iris.StatusInternalServerError, err500)
	app.OnErrorCode(iris.StatusNotFound, err404)

	app.Any("/ping", func(ctx iris.Context) {
		msg := "success"
		_, _ = ctx.WriteString(msg)
	})

	// /genids/getid
	genids := app.Party("/genids")
	genids.Get("/getid", api.IGenSnowId)

}
