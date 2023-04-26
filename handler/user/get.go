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

func Get(ctx *gin.Context) {
	log.Info("User Get function called")

	// TODO 改变user_id获取方式，或通过中间件实现
	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)
	place_id := user.Place_id

	// 解析请求数据
	begin_time_encode := ctx.Query("begin_time")
	end_time_encode := ctx.Query("end_time")

	var user_list []*model.User
	if begin_time_encode != "" && end_time_encode != "" {
		begin_time_decode, _ := url.QueryUnescape(begin_time_encode)
		end_time_decode, _ := url.QueryUnescape(end_time_encode)
		// 将字符串的时间转换为time.Time类型
		// TODO begin_time 和 end_time解析出来后都是UTC时间，但mysql中的时区未必是UTC，这里需要解决一下
		begin_time, _ := time.Parse(time.RFC3339, begin_time_decode)
		end_time, _ := time.Parse(time.RFC3339, end_time_decode)

		var err error
		user_list, err = model.FilterUserByTime(int(place_id), begin_time, end_time)
		if err != nil {
			// TODO 写入日志错误内容细化
			log.Error("filter user by time error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}
	} else {
		// 请求user_list数据
		var err error
		user_list, err = model.GetUserByPlaceId(int(place_id))
		if err != nil {
			// TODO 写入日志错误内容细化
			log.Error("search user_list_origin error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}
	}

	// 构造简短的用户信息
	var getResponse GetResponse
	for _, user := range user_list {

		// 管理员账户无需返回
		// if user_origin.Usertype == user_type.Admin {
		// 	continue
		// }

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
