package main

import (
	"meizuapi/thrift_gen/predict"
	"encoding/json"
	"github.com/spf13/cast"
	"time"
	"fmt"
	"crypto/md5"
	"encoding/hex"
)

var success int32= 1000
var SUCCESS  *int32 = &success

var mETA_DATA_ERROR int32= 1001
var META_DATA_ERROR *int32 = &mETA_DATA_ERROR

var uNKOWN_ERROR int32= 1002
var UNKOWN_ERROR *int32 = &uNKOWN_ERROR

var aUTHENTICATION_ERROR int32= 1003
var AUTHENTICATION_ERROR *int32 = &aUTHENTICATION_ERROR

var dATA_FORMAT_ERROR int32= 1004
var DATA_FORMAT_ERROR *int32 = &dATA_FORMAT_ERROR

var aD_NOT_FOUND int32= 110000
var AD_NOT_FOUND *int32 = &aD_NOT_FOUND

//FIMXME  注意sessionid、 的初始化

func initCallbackreq(thriftreq *predict.Req, req *CallbackReq){
	context_dict := make(map[string]string)
	impr_string_list := make([]string,0)
	/*
	注意指针类型
	 */
	thrift_ad_locations := make([] *predict.AdLocation,0)
	context_dict["net"] = req.Net
	context_dict["ip"] = req.Ip
	context_dict["device_model"] = req.DeviceModel
	context_dict["vc"] = req.VC
	context_dict["language"] = req.Language
	context_dict["uid"] = req.Uid
	context_dict["sn"] = req.Sn
	context_dict["request_id"] = req.ReqId

	//query := req.Query
	app_id := req.RealAppId
	var business int16
	if req.Business == nil {
		business = 1
	}else{
		business = *req.Business
	}

	is_idea := false
	for _,element := range req.Data{
		impr := make(map[string]interface{})
		thrift_ad_location := predict.NewAdLocation()
		thrift_ad_location.PositionId = element.PositionId
		thrift_ad_location.PositionId = element.PageId
		thrift_ad_locations = append(thrift_ad_locations,thrift_ad_location)
		impr["positon_id"] = element.PositionId
		impr["page_id"] = element.PageId
		impr["app_id"] = element.AppId
		impr["unit_id"] = element.UnitId
		impr["idea_id"] = element.IdeaId
		if impr["idea_id"] != 0 {
			is_idea = true
		}
		impr["score"] = element.Sorce
		impr["kw"] = element.Keyword
		impr_collection,_ := json.Marshal(impr)
		impr_string_list = append(impr_string_list,string(impr_collection))
 	}
	//FIXME adversion 如何使用  ??
 	//vtp := ""
 	//if app_id != 0{
 	//	vtp = "relate"
	//}else if is_idea{
	//	vtp = "idea"
	//}else if query != ""{
	//	vtp = "search"
	//}else{
	//	vtp = "display"
	//}
	//ab_version = ""
	source := 0
 	if is_idea{
 		source = 4
	}
	thriftreq.Business = cast.ToInt16(business)
	thriftreq.ReqId = &req.ReqId
	thriftreq.RelAppId = cast.ToInt64(app_id)
	thriftreq.Imei = req.Imei
	thriftreq.SessionId = req.SessionID
	thriftreq.AdLocations = thrift_ad_locations
	dictstring,_ := json.Marshal(context_dict)
	thriftreq.ContextInfo = string(dictstring)
	//thriftreq.AbtestParameters = 		//FIXME  ab test的更新
	thriftreq.Debug = true
	thriftreq.ImprList = impr_string_list
	thriftreq.Source = cast.ToInt16(source)
	//return thriftreq
}

func initPredictbinary(thriftreq *predict.Req, req *PredictReq){
	context_dict := make(map[string]string)
	thrift_ad_locations := make([] *predict.AdLocation,0)
	context_dict["net"] = req.Net
	context_dict["ip"] = req.Ip
	context_dict["device_model"] = req.DeviceModel
	context_dict["vc"] = req.VC
	context_dict["language"] = req.Language
	context_dict["uid"] = req.Uid
	context_dict["sn"] = req.Sn
	context_dict["request_id"] = req.ReqId

	var business int16
	var enableRank bool
	if req.EnableRank == nil{
		enableRank = true
	}else{
		enableRank = *req.EnableRank
	}
	if req.Business == nil {
		business = 1
	}else{
		business = *req.Business
	}
	for _,element := range req.AdLocations{
		thrift_ad_location := predict.NewAdLocation()
		thrift_ad_location.PositionId= element.PositionId
		thrift_ad_location.PageId = element.PageId
		thrift_ad_location.BlockId = element.BlockId
		thrift_ad_location.RankId = element.RankId
		thrift_ad_location.CategoryId = element.CategoryId
		thrift_ad_location.TagId = element.TagId
		thrift_ad_location.Position = element.Position
		thrift_ad_location.PositionType = element.PostionType
		thrift_ad_locations = append(thrift_ad_locations,thrift_ad_location)
	}
	thriftreq.ReqId = &req.ReqId
	thriftreq.Business = cast.ToInt16(business)
	thriftreq.Imei = req.Imei
	thriftreq.SessionId = req.SessionID
	thriftreq.AdLocations = thrift_ad_locations
	dictstring,_ := json.Marshal(context_dict)
	thriftreq.ContextInfo = string(dictstring)
	////thriftreq.AbtestParameters = 		//FIXME  ab test的更新
	thriftreq.EnableRank = enableRank
	thriftreq.Debug = true
	thriftreq.Topk = req.Topk
}

func initIdea_predict(thriftreq *predict.Req, req *Idea_predict){
	context_dict := make(map[string]string)
	thrift_ideas := make([] *predict.Idea,0)
	context_dict["net"] = req.Net
	context_dict["ip"] = req.Ip
	context_dict["device_model"] = req.DeviceModel
	context_dict["vc"] = req.VC
	context_dict["language"] = req.Language
	context_dict["uid"] = req.Uid
	context_dict["sn"] = req.Sn
	context_dict["request_id"] = req.ReqId
	for _,element := range req.Ideas{
		thrift_idea := predict.NewIdea()
		thrift_idea.IdeaId= element.IdeaId
		thrift_idea.Source = element.Source
		thrift_ideas = append(thrift_ideas,thrift_idea)
	}
	thriftreq.ReqId = &req.ReqId
	thriftreq.Imei = req.Imei
	thriftreq.SessionId = req.SessionID
	thriftreq.Ideas = thrift_ideas
	dictstring,_ := json.Marshal(context_dict)
	thriftreq.ContextInfo = string(dictstring)
	//thriftreq.AbtestParameters = 		//FIXME  ab test的更新
	thriftreq.Debug = true
	//
	//fmt.Print("\n   this is initIdea_predict")
	//fmt.Print(context_dict)
	//fmt.Print("\n this is  all ideas ",thriftreq.Ideas)

}

//var ab_config *TSMap
//func GetABTestParameters(version string) (string){
//	conf, exists := ab_config.Get(version)
//	if exists {
//		json := conf.(*json.RawMessage)
//		return fmt.Sprintf(`{"ab_version":"%s","3rd_rec":%s}`, version, string(*json))
//	}else{
//		return "{}"
//	}
//}
//
//func getABFromReq(c *gin.Context) string{
//	if ab_version,exits:= c.Get("version");exits{
//		return GetABTestParameters(cast.ToString(ab_version))
//	}
//	return "{}"
//}

func get_check_error_response(req *CheckRequst) (string, interface{}){
	data := map[string]interface{} {"reqid": cast.ToString(time.Now().Unix()), "code":META_DATA_ERROR, "msg":"failure"}
	request_id := req.ReqId
	appKey := req.Appkey
	checkinfo := req.Checkinfo
	nonce := req.Nonce
	if request_id !="" {
		data["reqid"] =  request_id
	}
	secret := "8c2f716071cea8233dec871fd7cefcaf"
	if request_id=="" || appKey=="" || checkinfo=="" || nonce==0 || appKey != "c7d805ab7f396ed6d420cd8da0f788ca"{
		return "", data
	}
	all_str_for_md5 :=cast.ToString(secret)+cast.ToString(appKey)+cast.ToString(nonce)
	h := md5.New()
	h.Write([]byte(all_str_for_md5))
	if hex.EncodeToString(h.Sum(nil)) != checkinfo{
		fmt.Print("\n 校验失败 \n")
		data["code"] = AUTHENTICATION_ERROR
		return request_id,  data
	}
	return request_id, nil
}

//参数再确定
//func get_predict_result(business int16, imei string, session_id string, ad_locations []map[interface{}]interface{}, context_dict map[string]interface{}, abtest_json_string string, enable_rank bool, is_debug bool, is_binary bool, topk int16) ([]byte,string){
//	var rank_results []byte   //is_binary
//	//var rank_results *toutiao_meizu_app.DisplayUnitList   //is_use_pb
//	var version string
//
//	//1.建立一个swrift链接
//	conn,_ := CallClient("binary_predict","","","")
//	//if conn != nil{
//	//	fmt.Print("conn is not nil \n\n")
//	//}
//	predict_server_client := getPredictClient(conn)
//
//	if predict_server_client != nil {
//		fmt.Print("predict_server_client != nil \n")
//	}
//
//	var thrift_ad_locations []*predict.AdLocation
//	for _, ad_location := range ad_locations {
//		thrift_ad_location := predict.NewAdLocation()
//		thrift_ad_location.PositionId = cast.ToInt64(ad_location["positionId"])
//		thrift_ad_location.PageId = cast.ToInt64(ad_location["pageId"])
//		//FIXME:  取零的用法？？
//		//thrift_ad_location.BlockID = ad_location.get("blockId", 0)
//		//thrift_ad_location.RankID = ad_location.get("rankId", 0)
//		//thrift_ad_location.CategoryID = ad_location.get("categoryId", 0)
//		//thrift_ad_location.TagID = ad_location.get("tagId", 0)
//		thrift_ad_location.BlockId = cast.ToInt64(ad_location["blockId"])
//		thrift_ad_location.RankId = cast.ToInt64(ad_location["rankId"])
//		thrift_ad_location.CategoryId = cast.ToInt64(ad_location["categoryId"])
//		thrift_ad_location.TagId = cast.ToInt64(ad_location["tagId"])
//		thrift_ad_location.Position = cast.ToInt32(ad_location["position"])
//		thrift_ad_location.PositionType = cast.ToInt32(ad_location["positionType"])
//		thrift_ad_locations = append(thrift_ad_locations, thrift_ad_location)
//	}
//	req := predict.NewReq()
//	req.Business = business
//	req.Imei = imei
//	req.SessionId = session_id
//	req.AdLocations = thrift_ad_locations
//	b, _ := json.Marshal(context_dict)
//	req.ContextInfo = string(b)
//	req.AbtestParameters = abtest_json_string
//	req.EnableRank = enable_rank
//	req.Debug = is_debug
//	req.Topk = topk
//
//	//2。 swrift 接口
//	//if is_binary {
//		rsp, exits := predict_server_client.BinaryPredict(req)
//	//}else{
//	//	rsp, _ := predict_server_client.Predict(defaultCtx, req)
//	//}
//	if exits != nil{
//		fmt.Print("\n")
//		fmt.Print("error")
//		fmt.Print("\n")
//	}
//
//	if rsp.Status == "success" {
//		if is_binary {
//			rank_results = rsp.RankResults //[]byte
//			//} else if is_use_pb { //proto 接口
//			//	rank_results = &toutiao_meizu_app.DisplayUnitList{}
//			//	for key, values := range rsp.RankResults {//FIXME  应该是Rsp类型的rank_results //这里格式如何转换？？
//			//		for _,value := range values {
//			//			rank_result := &toutiao_meizu_app.DisplayUnit{}
//			//			*(rank_result.PositionId) = cast.ToInt32(key)
//			//			*(rank_result.UnitId) = value.UnitID
//			//			*(rank_result.AppId) = value.AppId
//			//			if enable_rank {
//			//				*(rank_result.Ctr) = value.Score
//			//			} else {
//			//				*(rank_result.Ctr) = value.ctr
//			//			}
//			//			rank_results.Data = append(rank_results.Data,rank_result)
//			//			if is_truncated {
//			//				break
//			//			}
//			//		}
//			//		rank_results.Version = rsp.Version
//			//	}
//			//} else {//FIXME  应该是Rsp类型的rank_results //这里格式如何转换
//			//	for key, values := range rsp.RankResults {
//			//		unit_recommends := make([]interface{}, 0)
//			//		for value := range values {
//			//			unit_recommend := map[string]interface{}{"unitId": value.unit_id, "appId": value.app_id, "score": value.score, "price": value.price, "ctr": value.ctr, "relateScore": value.relevance}
//			//			unit_recommends = append(unit_recommends, unit_recommend)
//			//		}
//			//		rank_results[string(key)] = unit_recommends
//			//	}
//			//	if rsp.Vsersion == nil {
//			//		 version = ""
//			//	} else {
//			//		 version = rsp.Version
//			//	}
//			//}
//		}
//	}
//	return []byte(rank_results), version
//}


//参数再确定
//func callback_predict_result(callbackreq *Callbackreq,business int16, query string,imei string, app_id string,session_id string, ad_locations []map[string]interface{}, context_dict map[string]string, impr_string_list []string,abtest_json_string string,is_debug bool, source int16) (err error){
//	//1.建立一个swrift链接
//	conn,_ := CallClient("binary_predict","","","")
//	if conn != nil{
//		fmt.Print("conn != nil \n")
//	}
//	predict_server_client := getPredictClient(conn)
//
//	if predict_server_client != nil {
//		fmt.Print("predict_server_client != nil \n")
//	}
//	var thrift_ad_locations []*predict.AdLocation
//	for _, ad_location := range ad_locations {
//		thrift_ad_location := predict.NewAdLocation()
//		thrift_ad_location.PositionId = cast.ToInt64(ad_location["positionId"])
//		thrift_ad_location.PageId = cast.ToInt64(ad_location["pageId"])
//		thrift_ad_location.BlockId = cast.ToInt64(ad_location["blockId"])
//		thrift_ad_location.RankId = cast.ToInt64(ad_location["rankId"])
//		thrift_ad_location.CategoryId = cast.ToInt64(ad_location["categoryId"])
//		thrift_ad_location.TagId = cast.ToInt64(ad_location["tagId"])
//		thrift_ad_location.Position = cast.ToInt32(ad_location["position"])
//		thrift_ad_location.PositionType = cast.ToInt32(ad_location["positionType"])
//		thrift_ad_locations = append(thrift_ad_locations, thrift_ad_location)
//	}
//	req := predict.NewReq()
//	req.Business = business
//	fmt.Print("\n")
//	fmt.Print(reflect.TypeOf(context_dict["request_id"]))
//	fmt.Print("\n")
//	req.ReqId = &callbackreq.RequestId
//	fmt.Print("\n")
//	fmt.Print(*req.ReqId)
//	fmt.Print("\n")
//
//	req.RelAppId = cast.ToInt64(app_id)
//	req.Imei = imei
//	req.SessionId = session_id
//	req.AdLocations = thrift_ad_locations
//	b, _ := json.Marshal(context_dict)
//	req.ContextInfo = string(b)
//	req.AbtestParameters = abtest_json_string
//	req.ImprList = impr_string_list
//	req.Query = &callbackreq.Query
//	req.Debug = true
//	req.Source = source
//
//	//2。 swrift 接口
//	//if is_binary {
//	exits := predict_server_client.UploadServerImprOneway(req)
//	//}else{
//	//	rsp, _ := predict_server_client.Predict(defaultCtx, req)
//	//}
//	if exits != nil{
//		fmt.Print("\n")
//		fmt.Print("error")
//		fmt.Print("\n")
//	}
//	return nil
//}
//
