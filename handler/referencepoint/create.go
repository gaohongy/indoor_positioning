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

// TODO 由于添加rss时，可能因为没有对应的ap导致rss插入失败，但是参考点仍然会插入成功，是否考虑采用数据库的事务保证两个操作的原子性

// Create creates a new user account.
func Create(ctx *gin.Context) {
	log.Info("Referencepoint Create function called")

	var request CreateRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}
	// -------------------------------------------------------------------------------------
	// 创建参考点
	// TODO 改变user_id获取方式，或通过中间件实现
	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)

	place_id := user.Place_id

	// 创建参考点前，其所在的网格点不一定存在。查询所在网格点，当网格点不存在时，新建网格点
	gridpoint, err := model.GetGridpoint(request.Coordinate_x, request.Coordinate_y, request.Coordinate_z, place_id)
	// 查询结果为空err.Error() = "record not found"
	if err != nil {
		log.Error("gridpoint not exists, creating a new gridpoint", err)

		gridpoint = &model.Gridpoint{
			Coordinate_x: request.Coordinate_x,
			Coordinate_y: request.Coordinate_y,
			Coordinate_z: request.Coordinate_z,
			Place_id:     place_id,
			Createdate:   time.Now(),
			Updatedate:   time.Now(),
		}

		// TODO 网格点如果插入失败，这里直接return是否可以
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

	// TODO 验证参数合法性
	// if err := place.Validate(); err != nil {
	// 	handler.SendResponse(ctx, errno.ErrorValidation, nil)
	// 	return
	// }

	// 参考点数据插入数据库
	if err := referencepoint.Create(); err != nil {
		log.Error("referencepoint insert error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	referencepoint.Id = referencepoint.GetId()
	//------------------------------------------------------------------------------------------
	// 创建rss条目
	// TODO 打算用一个对象实现多次插入，但是在插入时会报主键重复的错误，或许相同对象插入时会看为相同条目
	// rss := &model.Rss{
	// 	Reference_point_id: referencepoint.Id,
	// 	Createdate:         time.Now(),
	// 	Updatedate:         time.Now(),
	// }

	// 依次处理列表中的数据
	for _, fingerPrint := range request.Rss_list {
		var ap_id uint64
		// 查询ap_id,如果ap不存在，则跳过，这里没办法自动添加ap，因为要确定ap位置，必须要人工添加
		if ap, err := model.GetAp(fingerPrint.Bssid); err != nil {
			continue
		} else {
			ap_id = ap.Id
		}

		rss := &model.Rss{
			Rss:                fingerPrint.Rss,
			Reference_point_id: referencepoint.Id,
			Ap_id:              ap_id,
			Createdate:         time.Now(),
			Updatedate:         time.Now(),
		}

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
