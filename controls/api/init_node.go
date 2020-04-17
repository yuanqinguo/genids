package api

import (
	"genids/config"
	"genids/idworker"
	"genids/utils"
	"github.com/kataras/iris/v12"
)

func InitNode(ctx iris.Context) {
	defer ctx.Next()
	extMsg := ""
	nodeid, err := ctx.URLParamInt64("nodeid")
	if err != nil || nodeid > 2 || nodeid < 0 {
		if err != nil {
			extMsg = err.Error()
		}
		SetResp(ctx, nil, utils.IDInvalidErr, extMsg)
		return
	}

	config.NodeID = nodeid
	idworker.GetIdWokrer()
	SetResp(ctx, "", utils.Ok, extMsg)
}
