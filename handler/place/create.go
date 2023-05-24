package place

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	Create
// @description	新建场所API
// @auth	高宏宇
// @param	ctx *gin.Context
func Create(ctx *gin.Context) {
	log.Info("Place Create function called")

	// 解析body参数
	var request CreateRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}

	place := model.Place{
		Place_address: request.Place_address,
		Longitude:     request.Longitude,
		Latitude:      request.Latitude,
		Map_id:        request.Map_id,
		Createdate:    time.Now(),
		Updatedate:    time.Now(),
	}

	// 场所数据插入数据库
	if err := place.Create(); err != nil {
		log.Error(errno.ErrorDatabase.Error(), err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	createResponse := CreateResponse{
		Place_id: place.GetId(),
	}

	// 发送响应
	handler.SendResponse(ctx, nil, createResponse)
}
