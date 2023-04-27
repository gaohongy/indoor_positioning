package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

func GetLocation(ctx *gin.Context) {
	log.Info("User Get Locatoin function called")

	// TODO 改变user_id获取方式，或通过中间件实现
	// content, _ := token.ParseRequest(ctx)
	// user, _ := model.GetUserById(content.ID)
	// place_id := user.Place_id

	// 解析请求数据
	user_id_str := ctx.Query("id")
	user_id_int, _ := strconv.Atoi(user_id_str)
	user_id := uint64(user_id_int)

	// 查询该用户最新的pathpoint
	pathpoint, err := model.FilterLatestPathpointByUserId(user_id)
	if err != nil {
		log.Error("FilterLatestPathpointByUserId error", err)
		handler.SendResponse(ctx, errno.ErrorRecordNotFound, nil)
		return
	}

	// 查询该pathpoint的详细坐标
	gridpoint, err := model.GetGridpointById(pathpoint.Grid_point_id)
	if err != nil {
		log.Error("FilterLatestPathpointByUserId error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 返回数据结构
	pathpoint_detail := model.Pathpoint_Detail{
		Id:           pathpoint.Id,
		Coordinate_x: gridpoint.Coordinate_x,
		Coordinate_y: gridpoint.Coordinate_y,
		Coordinate_z: gridpoint.Coordinate_z,
		Createdate:   pathpoint.Createdate,
	}

	handler.SendResponse(ctx, nil, pathpoint_detail)
}
