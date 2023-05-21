package referencepoint

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	Create
// @description	新建参考点API
// @auth	高宏宇
// @param	ctx *gin.Context
func Create(ctx *gin.Context) {
	log.Info("Referencepoint Create function called")

	// 解析body参数
	var request CreateRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}
	// -------------------------------------------------------------------------------------
	// 创建参考点
	// 获取登录用户ID
	content, _ := token.ParseRequest(ctx)
	// 查询用户
	user, _ := model.GetUserById(content.ID)
	// 查询用户所在场所ID
	place_id := user.Place_id

	// 创建参考点前，其所在的网格点不一定存在。查询所在网格点，当网格点不存在时，新建网格点
	gridpoint, err := model.GetGridpoint(request.Coordinate_x, request.Coordinate_y, request.Coordinate_z, place_id)
	// 查询结果为空err.Error() = "record not found"
	if err != nil {
		log.Info("gridpoint not exists, creating a new gridpoint")

		gridpoint = &model.Gridpoint{
			Coordinate_x: request.Coordinate_x,
			Coordinate_y: request.Coordinate_y,
			Coordinate_z: request.Coordinate_z,
			Place_id:     place_id,
			Createdate:   time.Now(),
			Updatedate:   time.Now(),
		}

		// 创建网格点
		if err := gridpoint.Create(); err != nil {
			log.Error("gridpoint insert error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}

		// 新插入网格点需要获取id
		gridpoint.Id = gridpoint.GetId()
	}

	referencepoint := model.Referencepoint{
		Grid_point_id: gridpoint.Id,
		Place_id:      place_id,
		Createdate:    time.Now(),
		Updatedate:    time.Now(),
	}

	//------------------------------------------------------------------------------------------
	// 创建rss条目
	// 依次处理列表中的数据
	flag := true // 是否为第一次插入参考点
	for _, fingerPrint := range request.Rss_list {
		var ap_id uint64
		// 查询ap_id,如果ap不存在，则跳过，这里没办法自动添加ap，因为要确定ap位置，必须要人工添加
		if ap, err := model.GetApByBssid(fingerPrint.Bssid); err != nil {
			continue
		} else {
			ap_id = ap.Id
		}

		// 若采集数据中无合法接入点，则不会创建参考点。flag确保参考点仅会创建一次
		if flag {
			// 参考点数据插入数据库
			if err := referencepoint.Create(); err != nil {
				log.Error("referencepoint insert error", err)
				handler.SendResponse(ctx, errno.ErrorDatabase, nil)
				return
			}

			referencepoint.Id = referencepoint.GetId()

			flag = false
		}

		rss := &model.Rss{
			Rss:                fingerPrint.Rss,
			Reference_point_id: referencepoint.Id,
			Ap_id:              ap_id,
			Createdate:         time.Now(),
			Updatedate:         time.Now(),
		}

		// 接入点插入数据库
		if err := rss.Create(); err != nil {
			log.Error("rss insert error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}
	}
	//---------------------------------------------------------------------------------------

	// 返回参考点id
	createResponse := CreateResponse{
		Referencepoint_id: referencepoint.Id,
	}

	// 发送响应
	handler.SendResponse(ctx, nil, createResponse)
}
