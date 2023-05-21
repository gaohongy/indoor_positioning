package location

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"
	"io"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	Create
// @description	请求坐标计算
// @auth	高宏宇
// @param	ctx *gin.Context
func Get(ctx *gin.Context) {
	log.Info("Getlocation function called")

	// 获取登录用户ID
	content, _ := token.ParseRequest(ctx)
	// 查询用户
	user, _ := model.GetUserById(content.ID)
	// 查询用户所在场所ID
	place_id := user.Place_id

	// 查询当前场所所有的AP
	ap_list, _ := model.GetApByPlaceId(place_id)
	// 创建AP的BSSID集合
	ap_bssid_set := make(map[string]struct{})
	for _, ap := range *ap_list {
		ap_bssid_set[ap.Bssid] = struct{}{}
	}

	//----------------------------------------------------------------------------
	// 解析请求参数
	request := ctx.Query("fingerprint")
	if request == "" {
		log.Error("parameter not found", errno.ErrorMissingParameter)
		handler.SendResponse(ctx, errno.ErrorMissingParameter, nil)
		return
	}

	// 将fingerprint字符串进行base64解码获取原始字节数据
	decoded, err := base64.StdEncoding.DecodeString(request)
	if err != nil {
		log.Error("Parameter parsing error", errno.ErrorParameterParsing)
		handler.SendResponse(ctx, errno.ErrorParameterParsing, nil)
		return
	}

	// 将字节数据反序列化为Go数据结构[]Fingerprint
	var fingerprints []Fingerprint
	json.Unmarshal(decoded, &fingerprints)

	//----------------------------------------------------------------------------
	var online_rss []float64        // 单一待定位点采集到的各个接入点信号强度列表
	var online_rss_list [][]float64 // 全部待定位点采集到的各个接入点信号强度列表合集
	var online_bssid_list []string  // 待定位点采集到的接入点的BSSID列表

	flag := false
	for _, fingerprint := range fingerprints {
		// 仅保留处在AP数据库中的AP数据
		if _, ok := ap_bssid_set[fingerprint.Bssid]; ok {
			flag = true
			online_rss = append(online_rss, fingerprint.Rss)
			online_bssid_list = append(online_bssid_list, fingerprint.Bssid)
		}
	}
	online_rss_list = append(online_rss_list, online_rss)
	// fmt.Println(online_rss_list)
	// fmt.Println(online_bssid_list)

	// 待定位点采集的AP同AP数据库无交集
	if !flag {
		handler.SendResponse(ctx, errno.ErrorAlgorithmCount, nil)
		return
	}

	// fmt.Println(online_rss_list)   // [[-8 -8]]
	// fmt.Println(online_bssid_list) // [00-00-00-00-00-00 11-11-11-11-11-11]
	// return

	//----------------------------------------------------------------------------
	// 查询当前场所所有的参考点
	referencepoint_list, _, err := model.ListReferencepointByPlaceid(place_id, 0, 0)
	if err != nil {
		log.Error("referencepoint list error", err)
		return
	}
	// fmt.Println(referencepoint_list)

	var offline_rss_list [][]float64      // 全部参考点采集的各个接入点的信号强度值
	var offline_location_list [][]float64 // 全部参考点的坐标，和offline_rss_list的顺序相同

	for _, referencepoint := range referencepoint_list {
		// ----------------------------------------------------------------------------
		// 创建当前参考点采集的接入点，BSSID到信号强度的哈希表
		rss_map := make(map[string]float64)

		// 查询到该参考点下所有rss条目
		rss_list, _, _ := model.ListRssByReferencepointid(referencepoint.Id, 0, 0)
		// fmt.Println(rss_list)
		// continue

		for _, rss := range rss_list {
			// 查询该rss条目所属AP
			ap, _ := model.GetApById(rss.Ap_id)
			// fmt.Println(ap)

			// 为了确保离线数据和在线数据的ap对应，采用字典进行存储
			rss_map[ap.Bssid] = rss.Rss
		}
		// return
		// ----------------------------------------------------------------------------

		// 在处理在线数据时，已经对待定位点采集的AP进行了筛选，此时online_bssid_list中的AP都存在于AP数据库
		var offline_rss []float64
		for _, online_bssid := range online_bssid_list {
			// 当参考点存在在线数据中的ap时，将该ap在该参考点的rss值加入数组，否则加入-100
			// fmt.Println(online_bssid)

			// 当前bssid对应AP被当前参考点采集到
			if value, ok := rss_map[online_bssid]; ok {
				offline_rss = append(offline_rss, value)
				// 当前bssid对应AP未被当前参考点采集到
				// RSS不同范围含义
				// 大于 -50 dBm: 强信号
				// [-50 dBm, -60 dBm]: 良好信号
				// [-60 dBm, -70 dBm]: 中等信号
				// [-70 dBm, -80 dBm]: 弱信号
				// 小于 -80 dBm: 微弱信号
			} else {
				offline_rss = append(offline_rss, -100)
			}
		}
		offline_rss_list = append(offline_rss_list, offline_rss)

		// 当前参考点所在网格点
		gridpoint, _ := model.GetGridpointById(referencepoint.Grid_point_id)
		// 构造离线坐标数据
		var offline_location []float64
		offline_location = append(offline_location, gridpoint.Coordinate_x)
		offline_location = append(offline_location, gridpoint.Coordinate_y)
		offline_location = append(offline_location, gridpoint.Coordinate_z)
		offline_location_list = append(offline_location_list, offline_location)

		// fmt.Println(offline_rss_list)
		// fmt.Println(offline_location_list)
		// return
	}
	// fmt.Println(offline_rss_list)
	// fmt.Println(offline_location_list)
	// online_rss_list(在线待预测数据)：e.g. [-57.02534192 -58.64624555 -60.49532547 -63.5991607  -52.55254457 -51.10042539]
	// online_bssid_list:
	// offline_rss_list(离线数据)：[ [-37.19167206 -57.68451847 -68.1915871  -70.40390708 -55.78694283 -63.78534359], [-34.06766605 -53.24717087 -57.5191879  -57.37401746 -57.46015148 -74.80795447] ]
	// offline_location_list(离线数据)：[ [10 10], [20 10] ]

	//----------------------------------------------------------------------------------------------------------------------------------------------------------
	// 坐标计算请求数据结构
	createKnnRequest := CreateKnnRequest{
		Offline_rss:      offline_rss_list,
		Offline_location: offline_location_list,
		Online_rss:       online_rss_list,
	}
	// fmt.Println(offline_rss_list)
	// fmt.Println(offline_location_list)
	// fmt.Println(online_rss_list)

	// 请求数据结构序列化为json数据
	jsonData, _ := json.Marshal(createKnnRequest)

	url := "http://localhost:8000/count"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error("make request python api error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 发起请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("request python api error", err)
		return
	}

	// 解析响应数据
	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		// 解析python api返回参数
		var createKnnResponse CreateKnnResponse
		json.Unmarshal(body, &createKnnResponse)

		// 向客户端返回数据
		var getResponse GetResponse
		for _, coordinate_list := range createKnnResponse.Coordinate {
			getResponse.Coordinate_x = math.Floor(coordinate_list[0])
			getResponse.Coordinate_y = math.Floor(coordinate_list[1])
			getResponse.Coordinate_z = math.Floor(coordinate_list[2])
		}

		handler.SendResponse(ctx, nil, getResponse)
		return
	}

	handler.SendResponse(ctx, errno.ErrorAlgorithmCount, nil)
}
