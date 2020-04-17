package mdw

import (
	"fmt"
	"genids/utils/logs"
	"github.com/getsentry/sentry-go"
	"github.com/kataras/iris/v12"
	"runtime"
	"time"
)

// 统一异常处理
func NewRecoverMdw() iris.Handler {
	return func(ctx iris.Context) {
		defer func() {
			if err := recover(); err != nil {
				sentry.CurrentHub().Recover(err)
				sentry.Flush(time.Second * 3)
				if ctx.IsStopped() {
					return
				}

				var stacktrace string
				for i := 1; ; i++ {
					_, f, l, got := runtime.Caller(i)
					if !got {
						break
					}

					stacktrace += fmt.Sprintf("%s:%d\n", f, l)
				}

				// when stack finishes
				logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.HandlerName())
				logMessage += fmt.Sprintf("At Request: %d %s %s %s\n", ctx.GetStatusCode(), ctx.Path(), ctx.Method(), ctx.RemoteAddr())
				logMessage += fmt.Sprintf("Trace: %s\n", err)
				logMessage += fmt.Sprintf("\n%s", stacktrace)

				logs.LogSystem.Errorf("recover => %s", logMessage)

				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.StopExecution()
			}
		}()

		ctx.Next()
	}
}

func NewResponseMdw() iris.Handler {
	return func(ctx iris.Context) {
		res := ctx.Values().Get("resp")
		if res != nil {
			_, err := ctx.JSON(res)
			if err != nil {
				logs.LogSystem.Error("WriterResp error: ", err.Error())
			}
		}
		ctx.Next()
	}
}
