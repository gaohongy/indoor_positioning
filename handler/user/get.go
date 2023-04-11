package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"
	"indoor_positioning/pkg/user_type"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

func Get(ctx *gin.Context) {
	log.Info("User Get function called")

	// TODO 改变user_id获取方式，或通过中间件实现
	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)

	place_id := user.Place_id

	// 请求user_list数据
	user_list_origin, err := model.GetUserByPlaceId(int(place_id))
	if err != nil {
		// TODO 写入日志错误内容细化
		log.Error("search user_list_origin error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 构造简短的用户信息
	var user_brief_list []model.User_Brief
	for _, user_origin := range *user_list_origin {

		// 管理员账户无需返回
		if user_origin.Usertype == user_type.Admin {
			continue
		}

		user_brief := model.User_Brief{
			Id:       user_origin.Id,
			Username: user_origin.Username,
		}
		user_brief_list = append(user_brief_list, user_brief)
	}
	handler.SendResponse(ctx, nil, user_brief_list)
}
