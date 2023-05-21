package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	Get
// @description	查询符合时间条件的用户API
// @auth	高宏宇
// @param	ctx *gin.Context
func Get(ctx *gin.Context) {
	log.Info("User Get function called")

	// 获取登录用户ID
	content, _ := token.ParseRequest(ctx)
	// 查询用户
	user, _ := model.GetUserById(content.ID)
	// 查询用户所在场所ID
	place_id := user.Place_id

	// 解析请求参数
	begin_time_encode := ctx.Query("begin_time")
	end_time_encode := ctx.Query("end_time")

	var user_list []*model.User
	// 有时间筛选条件
	if begin_time_encode != "" && end_time_encode != "" {
		// 将url安全格式编码字符串还原为原始字符串
		begin_time_decode, _ := url.QueryUnescape(begin_time_encode)
		end_time_decode, _ := url.QueryUnescape(end_time_encode)
		// 将字符串的时间转换为time.Time类型
		begin_time, _ := time.Parse(time.RFC3339, begin_time_decode)
		end_time, _ := time.Parse(time.RFC3339, end_time_decode)

		// 带筛选条件查询用户
		var err error
		user_list, err = model.FilterUserByTime(int(place_id), begin_time, end_time)
		if err != nil {
			log.Error("filter user by time error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}
		// 无时间筛选条件
	} else {
		// 查询当前场所全部用户
		var err error
		user_list, err = model.GetUserByPlaceId(int(place_id))
		if err != nil {
			log.Error("search user_list_origin error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}
	}

	// 构造简短的用户信息
	var getResponse GetResponse
	for _, user := range user_list {

		user_brief := model.User_Brief{
			Id:         user.Id,
			Username:   user.Username,
			Usertype:   user.Usertype,
			Place_id:   user.Place_id,
			Createdate: user.Createdate,
			Updatedate: user.Updatedate,
		}
		getResponse.User_list = append(getResponse.User_list, user_brief)
	}
	handler.SendResponse(ctx, nil, getResponse)
}
