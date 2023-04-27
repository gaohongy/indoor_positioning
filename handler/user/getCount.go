package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"
	"net/url"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

func GetCount(ctx *gin.Context) {
	log.Info("User GetCount function called")

	// TODO 改变user_id获取方式，或通过中间件实现
	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)
	place_id := user.Place_id

	// 解析请求数据
	begin_time_encode := ctx.Query("begin_time")
	end_time_encode := ctx.Query("end_time")
	// unit := ctx.Query("unit") // TODO 添加对不同时间间隔单位的支持，目前是以天为间隔单位

	begin_time_decode, _ := url.QueryUnescape(begin_time_encode)
	end_time_decode, _ := url.QueryUnescape(end_time_encode)
	begin_time, _ := time.Parse(time.RFC3339, begin_time_decode)
	end_time, _ := time.Parse(time.RFC3339, end_time_decode)

	// 请求user_list数据
	user_list, err := model.FilterUserByTime(int(place_id), begin_time, end_time)
	if err != nil {
		// TODO 写入日志错误内容细化
		log.Error("filter user_list error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// date_list 存储map中的key值，方便排序
	// date_amount_map 存储时间和人数之间的映射
	var date_list []string
	date_amount_map := make(map[string]UserAmount)
	for _, user := range user_list {
		date_string := user.Createdate.Format("2006/01/02")
		_, ok := date_amount_map[date_string]
		if ok {
			useramount := date_amount_map[date_string]
			if user.Usertype == 0 { // 管理员
				date_amount_map[date_string] = UserAmount{AdminAmount: useramount.AdminAmount + 1, OrdinaryUserAmount: useramount.OrdinaryUserAmount, SumUserAmount: useramount.SumUserAmount + 1}
			} else {
				date_amount_map[date_string] = UserAmount{AdminAmount: useramount.AdminAmount, OrdinaryUserAmount: useramount.OrdinaryUserAmount + 1, SumUserAmount: useramount.SumUserAmount + 1}
			}
		} else {
			// 日期以前未出现，说明该日期对应的AdminAmount和OrdinaryUserAmount均不存在，直接赋值为1
			if user.Usertype == 0 {
				date_amount_map[date_string] = UserAmount{AdminAmount: 1, OrdinaryUserAmount: 0, SumUserAmount: 1}
			} else {
				date_amount_map[date_string] = UserAmount{AdminAmount: 0, OrdinaryUserAmount: 1, SumUserAmount: 1}
			}
			date_list = append(date_list, date_string) // 避免时间重复
		}
	}
	// 自定义排序方式
	less := func(i, j int) bool {
		date1, _ := time.ParseInLocation("2006/01/02", date_list[i], time.Local)
		date2, _ := time.ParseInLocation("2006/01/02", date_list[j], time.Local)
		return date1.Before(date2)
	}
	// 将时间字符串按照时间先后顺序排序
	sort.Slice(date_list, less)

	// 根据排序后的顺序构造返回数据，确保返回数据中的先后顺序符合时间先后顺序
	var user_amount_list []GetUserAmountResponse
	for _, date_string := range date_list {
		date, _ := time.ParseInLocation("2006/01/02", date_string, time.Local)
		useramount := date_amount_map[date_string]
		user_amount := GetUserAmountResponse{
			Date:               date,
			AdminAmount:        useramount.AdminAmount,
			OrdinaryUserAmount: useramount.OrdinaryUserAmount,
			SumUserAmount:      useramount.SumUserAmount,
		}
		user_amount_list = append(user_amount_list, user_amount)
	}

	handler.SendResponse(ctx, nil, user_amount_list)
}
