package place

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	GetMap
// @description	查询场所对应地图API
// @auth	高宏宇
// @param	ctx *gin.Context
func GetMap(ctx *gin.Context) {
	log.Info("place GetMap function called")

	// 获取登录用户ID
	content, _ := token.ParseRequest(ctx)
	// 查询用户
	user, _ := model.GetUserById(content.ID)
	// 查询用户所在场所ID
	place_id := user.Place_id

	// 查询场所
	place, err := model.GetPlaceById(place_id)
	if err != nil {
		log.Error("search place error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 响应数据
	getMapIdResponse := GetMapIdResponse{
		Map_id: place.Map_id,
	}

	handler.SendResponse(ctx, nil, getMapIdResponse)
}
