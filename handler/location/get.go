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

// Create creates a new user account.
func Get(ctx *gin.Context) {
	log.Info("Getlocation function called")

	// TODO 改变user_id获取方式，或通过中间件实现
	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)
	place_id := user.Place_id

	// 查询当前场所所有的AP，创建AP的BSSID集合
	ap_list, _ := model.GetApByPlaceId(place_id)
	ap_bssid_set := make(map[string]struct{})
	for _, ap := range *ap_list {
		ap_bssid_set[ap.Bssid] = struct{}{}
	}

	//----------------------------------------------------------------------------
	// 解析请求参数
	request := ctx.Query("fingerprint")
	// TODO 缺少参数时是不是空值
	if request == "" {
		log.Error("parameter not found", errno.ErrorMissingParameter)
		handler.SendResponse(ctx, errno.ErrorMissingParameter, nil)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(request)
	if err != nil {
		log.Error("Parameter parsing error", errno.ErrorParameterParsing)
		handler.SendResponse(ctx, errno.ErrorParameterParsing, nil)
		return
	}

	// json数组解析为[]Fingerprint
	var fingerprints []Fingerprint
	json.Unmarshal(decoded, &fingerprints)

	//----------------------------------------------------------------------------
	var online_rss []float64 // 在线rss数据
	var online_rss_list [][]float64
	var online_bssid_list []string
	flag := false
	for _, fingerprint := range fingerprints {
		if _, ok := ap_bssid_set[fingerprint.Bssid]; ok { // 仅保留处在AP数据库中的AP数据
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

	// TEST 在线数据
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

	// TEST referencepoint_list
	// fmt.Println(referencepoint_list)
	// return

	var offline_rss_list [][]float64      // 离线rss数据，多个参考点的，和在线数据中ap对应的rss值
	var offline_location_list [][]float64 // 和上述rss数组对应的坐标

	for _, referencepoint := range referencepoint_list {
		// 当前参考点所在网格点
		gridpoint, _ := model.GetGridpointById(referencepoint.Grid_point_id)

		rss_map := make(map[string]float64)

		// 查询到该参考点下所有rss条目
		rss_list, _, _ := model.ListRssByReferencepointid(referencepoint.Id, 0, 0)
		// TEST rss_list
		// fmt.Println(rss_list)
		// continue

		for _, rss := range rss_list {
			// 查询该rss条目所属AP
			ap, _ := model.GetApById(rss.Ap_id)
			// TEST ap
			// fmt.Println(ap)

			// 为了确保离线数据和在线数据的ap对应，采用字典进行存储
			rss_map[ap.Bssid] = rss.Rss
		}
		// return

		// 在处理在线数据时，已经对待定位点采集的AP进行了筛选，此时online_bssid_list中的AP都存在于AP数据库
		var offline_rss []float64
		for _, online_bssid := range online_bssid_list {
			// 当参考点存在在线数据中的ap时，将该ap在该参考点的rss值加入数组，否则加入0
			// fmt.Println(online_bssid)
			if value, ok := rss_map[online_bssid]; ok { // 当前bssid对应AP被当前参考点采集到
				offline_rss = append(offline_rss, value)
			} else {
				offline_rss = append(offline_rss, 0)
			}
		}

		offline_rss_list = append(offline_rss_list, offline_rss)

		var offline_location []float64
		offline_location = append(offline_location, gridpoint.Coordinate_x)
		offline_location = append(offline_location, gridpoint.Coordinate_y)
		offline_location = append(offline_location, gridpoint.Coordinate_z)
		offline_location_list = append(offline_location_list, offline_location)

		// TEST offline_rss_list offline_location_list
		// fmt.Println(offline_rss_list)
		// fmt.Println(offline_location_list)
		// return
	}
	// TEST offline_rss_list offline_location_list
	// fmt.Println(offline_rss_list)
	// fmt.Println(offline_location_list)
	// return
	// 截止到目前为止，我们得到了如下数据：
	// online_rss_list(在线待预测数据)：e.g. [-57.02534192 -58.64624555 -60.49532547 -63.5991607  -52.55254457 -51.10042539]
	// online_bssid_list:
	// offline_rss_list(离线数据)：[ [-37.19167206 -57.68451847 -68.1915871  -70.40390708 -55.78694283 -63.78534359], [-34.06766605 -53.24717087 -57.5191879  -57.37401746 -57.46015148 -74.80795447] ]
	// offline_location_list(离线数据)：[ [10 10], [20 10] ]

	//----------------------------------------------------------------------------------------------------------------------------------------------------------
	// TODO 把数据传递给python 因为只是内部传递调用，因此get方法采用body参数也无妨(其他方式太难传递参数)
	// //
	createKnnRequest := CreateKnnRequest{
		Offline_rss:      offline_rss_list,
		Offline_location: offline_location_list,
		Online_rss:       online_rss_list,
	}
	// TEST
	// fmt.Println(offline_rss_list)
	// fmt.Println(offline_location_list)
	// fmt.Println(online_rss_list)

	jsonData, _ := json.Marshal(createKnnRequest)
	print(jsonData)
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
		// TODO 这里后续要考虑处理多组knn计算结果的返回，目前认为每次只是单点计算
		// NOTE gridsize的存在说明网格密度是固定的，即基本单位是0.01m，目前把所有数据均扩大100倍，也就是存储的数据不存在小数，knn的计算结果出现小数是因为权重的存在，但是实际存储只存储整数部分
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

// req, _ := http.NewRequest(http.MethodGet, "localhost/count", nil)

// // 设置url参数
// params := make(url.Values)
// params.Add("offline_rss_list", string(offline_rss_list_json))
// params.Add("offline_location_list", string(offline_location_list_json))
// params.Add("online_rss_list", string(online_rss_list_json))
// req.URL.RawQuery = params.Encode()

// TODO 需要考虑在接口中加入knn中k值的设定
