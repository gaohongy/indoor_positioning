package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	GetLocation
// @description	查询用户最新位置API
// @auth	高宏宇
// @param	ctx *gin.Context
func GetLocation(ctx *gin.Context) {
	log.Info("User Get Locatoin function called")

	// 解析请求数据
	user_id_str := ctx.Query("id")              // user_id_str: 字符串类型
	user_id_int, _ := strconv.Atoi(user_id_str) // user_id_int: int类型
	user_id := uint64(user_id_int)              // user_id: uint64类型

	// 查询该用户最新的路径点
	pathpoint, err := model.FilterLatestPathpointByUserId(user_id)
	if err != nil {
		log.Error("FilterLatestPathpointByUserId error", err)
		handler.SendResponse(ctx, errno.ErrorRecordNotFound, nil)
		return
	}

	// 查询该路径点所在网格点
	gridpoint, err := model.GetGridpointById(pathpoint.Grid_point_id)
	if err != nil {
		log.Error("FilterLatestPathpointByUserId error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 构造响应数据结构
	pathpoint_detail := model.Pathpoint_Detail{
		Id:           pathpoint.Id,
		Coordinate_x: gridpoint.Coordinate_x,
		Coordinate_y: gridpoint.Coordinate_y,
		Coordinate_z: gridpoint.Coordinate_z,
		Createdate:   pathpoint.Createdate,
	}

	handler.SendResponse(ctx, nil, pathpoint_detail)
}
