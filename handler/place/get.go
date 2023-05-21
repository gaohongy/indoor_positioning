package place

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	Get
// @description	查询全部场所API
// @auth	高宏宇
// @param	ctx *gin.Context
func Get(ctx *gin.Context) {
	log.Info("place Get function called")

	// 查询当前场所所有的参考点
	place_list_origin, err := model.GetAllPlaces()
	if err != nil {
		log.Error("search places list error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 构造简洁的place_list
	var place_brief_list []model.Place_brief
	for _, place_origin := range place_list_origin {

		place_brief := model.Place_brief{
			Id:            place_origin.Id,
			Place_address: place_origin.Place_address,
			Longitude:     place_origin.Longitude,
			Latitude:      place_origin.Latitude,
		}
		place_brief_list = append(place_brief_list, place_brief)
	}

	handler.SendResponse(ctx, nil, place_brief_list)
}
