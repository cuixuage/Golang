package main

import (
	"meizuapi/thrift_gen/predict"
	"encoding/json"
	"github.com/spf13/cast"
	"time"
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"code.byted.org/gopkg/thrift"
	"net"
	"2018_2_5/abtest_thrift"
)
var success int32= 1000
var SUCCESS  *int32 = &success
var StrSuccess = "success"

var META_DATA_ERROR int32= 1001

var uNKOWN_ERROR int32= 1002
var UNKOWN_ERROR *int32 = &uNKOWN_ERROR

var aUTHENTICATION_ERROR int32= 1003
var AUTHENTICATION_ERROR *int32 = &aUTHENTICATION_ERROR

var dATA_FORMAT_ERROR int32= 1004
var DATA_FORMAT_ERROR *int32 = &dATA_FORMAT_ERROR

var aD_NOT_FOUND int32= 110000
var AD_NOT_FOUND *int32 = &aD_NOT_FOUND

const MEIZU_UID_TYPE = 111

//func _c_mul(a, b) {
//	return eval(hex((long(a) * b) & (2 * *64 - 1))[:-1])
//}

//func hashstring(my_string) {
//	if not my_string:
//	return 0 # empty
//	value = ord(my_string[0]) << 7
//	for
//	char
//	in
//my_string:
//	value = _c_mul(1000003, value) ^ ord(char)
//	value = value ^ len(my_string)
//	return value % (2 * * 64)
//}

func getABClient(conn net.Conn) *abtest_thrift.VersionServiceClient {
	var transport thrift.TTransport
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTBufferedTransportFactory(4096))
	transport = thrift.NewTSocketFromConnTimeout(conn, TIMEOUT_CONNECT_TIME * time.Millisecond)
	transport = transportFactory.GetTransport(transport)
	return abtest_thrift.NewVersionServiceClientFactory(transport, protocolFactory)
}

func get_ab_setting(){
	//request = {"token":"3rd_meizu_app","uid":hashstring(imei), "uid_type":MEIZU_UID_TYPE}
	request := `{"token":"3rd_meizu_app","uid":"", "uid_type":MEIZU_UID_TYPE}`
	thriftReq := &abtest_thrift.VersionReq{request}
	//localIP := kite.GetLocalIp()
	localIP := "10.8.64.231"
	fmt.Print(localIP,"\n")
	port := "7701"
	c,err := pool.Get(localIP, port, TIMEOUT_CONNECT_TIME*time.Millisecond)
	if err != nil{
		fmt.Print(err," ","can't get conn \n")
	}
	abClient := getABClient(c)
	rsp,err := abClient.GetVersionWithUserRequestForGolang(thriftReq)
	if err != nil{
		fmt.Print("can't get ab_response\n")
	}
	fmt.Print(rsp.Err," ",rsp.Info," ",*rsp.Msg,"\n");
}

func initCallbackreq(thriftreq *predict.Req, req *CallbackReq){
	context_dict := make(map[string]interface{})
	impr_string_list := make([]string,0)
	thrift_ad_locations := make([] *predict.AdLocation,0)
	context_dict["net"] = req.Net
	context_dict["ip"] = req.Ip
	if req.DeviceModel != "" {
		context_dict["device_model"] = req.DeviceModel
	}
	if req.VC != ""{
		context_dict["vc"] = req.VC
	}
	if req.Language != ""{
		context_dict["language"] = req.Language
	}
	if req.Uid != ""{
		context_dict["uid"]= req.Uid
	}
	if req.Sn != ""{
		context_dict["sn"] = req.Sn
	}

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
		impr["score"] = element.Score
		impr["kw"] = element.Keyword
		impr_collection,_ := json.Marshal(impr)
		impr_string_list = append(impr_string_list,string(impr_collection))
 	}

 	//fmt.Print("\nAB test func \n")
	//get_ab_setting()


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
	thriftreq.Business = business
	thriftreq.ReqId = &req.ReqId
	thriftreq.RelAppId = app_id
	thriftreq.Imei = req.Imei
	thriftreq.SessionId = req.SessionID
	thriftreq.AdLocations = thrift_ad_locations
	dictstring,_ := json.Marshal(context_dict)
	thriftreq.ContextInfo = string(dictstring)
	//thriftreq.AbtestParameters = 		//FIXME  ab test的更新
	thriftreq.ImprList = impr_string_list
	thriftreq.Source = cast.ToInt16(source)
	//return thriftreq
}

func initPredictbinary(thriftreq *predict.Req, req *PredictReq){
	context_dict := make(map[string]interface{})
	thrift_ad_locations := make([] *predict.AdLocation,0)
	context_dict["net"] = req.Net
	context_dict["ip"] = req.Ip
	//if req.DeviceModel != "" {
	//	context_dict["device_model"] = req.DeviceModel
	//}
	if req.VC != 0{
		context_dict["vc"] = req.VC
	}
	if req.Language != ""{
		context_dict["language"] = req.Language
	}
	if req.Uid != ""{
		context_dict["uid"]= req.Uid
	}
	if req.Sn != ""{
		context_dict["sn"] = req.Sn
	}

	var business int16
	var enableRank bool
	var topk int16
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
	if req.Topk == nil {
		topk = 0
	}else{
		topk = *req.Topk
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

	thriftreq.Business = business
	thriftreq.Imei = req.Imei
	thriftreq.SessionId = req.SessionID
	thriftreq.AdLocations = thrift_ad_locations
	dictstring,_ := json.Marshal(context_dict)
	thriftreq.ContextInfo = string(dictstring)
	//thriftreq.AbtestParameters = "" 		//FIXME  ab test的更新
	thriftreq.EnableRank = enableRank
	thriftreq.Topk = topk
}

func initRecommend(thriftreq *predict.Req, req *RecommendReq){
	fmt.Print("init thrift nil\n")
	fmt.Print("init thrift nil\n")
	fmt.Print("init thrift nil\n")


	context_dict := make(map[string]interface{})
	thrift_ad_locations := make([] *predict.AdLocation,0)
	context_dict["net"] = req.Net
	context_dict["ip"] = req.Ip
	//if req.DeviceModel != "" {
	//	context_dict["device_model"] = req.DeviceModel
	//}
	if req.VC != 0{
		context_dict["vc"] = req.VC
	}
	if req.Language != ""{
		context_dict["language"] = req.Language
	}
	if req.Uid != ""{
		context_dict["uid"]= req.Uid
	}
	if req.Sn != ""{
		context_dict["sn"] = req.Sn
	}

	var business int16
	var enableRank bool
	var topk int16
	var appid int64
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
	if req.Topk == nil {
		topk = 0
	}else{
		topk = *req.Topk
	}
	if req.AppId == nil {
		appid = 0
	}else{
		appid = *req.AppId
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

	thriftreq.Business = business
	thriftreq.Imei = req.Imei
	thriftreq.SessionId = req.SessionID
	thriftreq.AdLocations = thrift_ad_locations
	dictstring,_ := json.Marshal(context_dict)
	thriftreq.ContextInfo = string(dictstring)
	//thriftreq.AbtestParameters = "" 		//FIXME  ab test的更新
	thriftreq.EnableRank = enableRank
	thriftreq.Topk = topk
	thriftreq.RelAppId = appid

}
func initIdea_predict(thriftreq *predict.Req, req *Idea_predict){
	context_dict := make(map[string]interface{})
	thrift_ideas := make([] *predict.Idea,0)
	context_dict["net"] = req.Net
	context_dict["ip"] = req.Ip
	//if req.DeviceModel != "" {
	//	context_dict["device_model"] = req.DeviceModel
	//}
	//if req.VC != ""{
	//	context_dict["vc"] = req.VC
	//}
	if req.Language != ""{
		context_dict["language"] = req.Language
	}
	if req.Uid != ""{
		context_dict["uid"]= req.Uid
	}
	if req.Sn != ""{
		context_dict["sn"] = req.Sn
	}

	for _,element := range req.Ideas{
		thrift_idea := predict.NewIdea()
		thrift_idea.IdeaId= element.IdeaId
		thrift_idea.Source = element.Source
		thrift_ideas = append(thrift_ideas,thrift_idea)
	}
	thriftreq.Imei = req.Imei
	thriftreq.SessionId = req.SessionID
	thriftreq.Ideas = thrift_ideas
	dictstring,_ := json.Marshal(context_dict)
	thriftreq.ContextInfo = string(dictstring)
	//thriftreq.AbtestParameters = 		//FIXME  ab test的更新
	//thriftreq.Debug = true
	//
}

func initSearch (thriftreq *predict.Req, req *Search){
	context_dict := make(map[string]interface{})
	context_dict["net"] = req.Net
	context_dict["ip"] = req.Ip
	//if req.DeviceModel != "" {
	//	context_dict["device_model"] = req.DeviceModel
	//}
	if req.VC != 0{
		context_dict["vc"] = req.VC
	}
	if req.Language != ""{
		context_dict["language"] = req.Language
	}
	if req.Uid != 0{
		context_dict["uid"]= req.Uid
	}
	if req.Sn != ""{
		context_dict["sn"] = req.Sn
	}

	var business int16
	var enableRank bool
	var source int16
	keyword := req.Kw        //*string类型
	if req.EnableRank == nil{
		enableRank = true
	}else{
		enableRank = *req.EnableRank
	}
	if req.Source == nil {
		source = 1
	}else{
		source = *req.Source
	}
	if req.Business == nil {
		business = 1
	}else{
		business = *req.Business
	}

	thriftreq.Business = business
	thriftreq.Imei = req.Imei
	thriftreq.SessionId = req.SessionID
	thriftreq.Query = keyword
	search := "search"
	thriftreq.QueryType = &search
	thriftreq.AdLocations = make([]*predict.AdLocation,0)
	dictstring,_ := json.Marshal(context_dict)
	thriftreq.ContextInfo = string(dictstring)
	//thriftreq.AbtestParameters = 		//FIXME  ab test的更新
	thriftreq.EnableRank = enableRank
	thriftreq.Source = source
}


func initUpload_behavior (thriftreq *predict.Req, req *Search){
	context_dict := make(map[string]interface{})
	context_dict["net"] = req.Net
	context_dict["ip"] = req.Ip
	//if req.DeviceModel != "" {
	//	context_dict["device_model"] = req.DeviceModel
	//}
	if req.VC != 0{
		context_dict["vc"] = req.VC
	}
	if req.Language != ""{
		context_dict["language"] = req.Language
	}
	if req.Uid != 0{
		context_dict["uid"]= req.Uid
	}
	if req.Sn != ""{
		context_dict["sn"] = req.Sn
	}

	var business int16
	var enableRank bool
	var source int16
	keyword := req.Kw        //*string类型
	if req.EnableRank == nil{
		enableRank = true
	}else{
		enableRank = *req.EnableRank
	}
	if req.Source == nil {
		source = 1
	}else{
		source = *req.Source
	}
	if req.Business == nil {
		business = 1
	}else{
		business = *req.Business
	}

	thriftreq.Business = business
	thriftreq.Imei = req.Imei
	thriftreq.SessionId = req.SessionID
	thriftreq.Query = keyword
	search := "search"
	thriftreq.QueryType = &search
	thriftreq.AdLocations = make([]*predict.AdLocation,0)
	dictstring,_ := json.Marshal(context_dict)
	thriftreq.ContextInfo = string(dictstring)
	//thriftreq.AbtestParameters = 		//FIXME  ab test的更新
	thriftreq.EnableRank = enableRank
	thriftreq.Source = source
}

func MD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func get_check_error_response(req *CheckRequst) (string, interface{}){
	data := map[string]interface{} {"reqid": cast.ToString(time.Now().Unix()), "code":META_DATA_ERROR, "msg":"failure"}
	request_id := req.ReqId
	appKey := req.Appkey
	checkinfo := req.Checkinfo
	nonce := req.Nonce
	//fmt.Print(" Metadata==")
	//fmt.Print(" ",request_id," ",appKey," ",checkinfo," ",nonce,"\n")
	if request_id !="" {
		data["reqid"] =  request_id
	}
	secret := "8c2f716071cea8233dec871fd7cefcaf"
	if request_id == "" || checkinfo == "" || nonce == "" || nonce == 0 || appKey != "c7d805ab7f396ed6d420cd8da0f788ca" {
		return "", data
	}
	all_str_for_md5 :=cast.ToString(secret)+cast.ToString(appKey)+cast.ToString(nonce)
	h := md5.New()
	h.Write([]byte(all_str_for_md5))
	//fmt.Print(MD5(all_str_for_md5)," === ",checkinfo," === ",hex.EncodeToString(h.Sum(nil)))
	if hex.EncodeToString(h.Sum(nil)) != checkinfo{
		//fmt.Print("checkHeader error \n")
		data["code"] = AUTHENTICATION_ERROR
		return request_id,  data
	}
	if MD5(all_str_for_md5)==checkinfo && checkinfo==hex.EncodeToString(h.Sum(nil)){
	  //fmt.Print("checkHeader success\n")
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
