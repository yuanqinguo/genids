package api

import (
	"genids/idworker"
	"genids/utils"
	"github.com/kataras/iris/v12"
)

func IGenSnowId(ctx iris.Context) {
	defer ctx.Next()
	var idGen *idworker.IdWorker
	for i := 0; i < 5; i++ {
		idGen = idworker.GetIdWokrer()
		if idGen != nil {
			break
		}
	}

	if idGen == nil {
		SetResp(ctx, "", utils.IDWorkErr, "")
		return
	}
	id, err := idGen.NextId()
	if err != nil {
		SetResp(ctx, "", utils.NextIdErr, err.Error())
		return
	}
	SetResp(ctx, struct {
		Id int64 `json:"id"`
	}{Id: id}, utils.Ok, "")
}
