package api

import (
	"fmt"
	"genids/utils"
	"github.com/kataras/iris/v12"
)

type Resp struct {
	Code   int64       `json:"errcode"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

//在这里统一封装回复
func SetResp(ctx iris.Context, data interface{}, errno *utils.Errno, extMsg string) {
	formatExtMsg := func(extMsg string) string {
		if len(extMsg) > 0 {
			return fmt.Sprintf("[%s]", extMsg)
		}
		return ""
	}
	ctx.Values().Set("resp", Resp{
		Code:   errno.Code,
		Errmsg: errno.Message + formatExtMsg(extMsg),
		Data:   data,
	})
}
