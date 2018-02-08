package main

import (
	"github.com/gin-gonic/gin"
	"meizuapi/thrift_gen/predict"
	"fmt"
)
/*
格式标注还挺严格的  json
后面判断是否为空的参数  转换为指针再判断nil
 */
type CallbackReq struct{
	SessionID string 	`json:"session_id"`
	ReqId  string 		`json:"requestId"`
	Imei   string 		`json:"imei"`
	Query  string 		`json:"query"`
	Data []struct{
		Sorce float64		`json:"score"`
		IdeaId int32		`json:"idea_id"`
		PositionId int64 	`jso:"position_id"`
		PageId int64 		`json:"page_id"`
		AppId int64			`json:"app_id"`
		UnitId int64		`json:"unit_id"`
		Keyword string 		`json:"kw"`
 	} `json:"data"`
	DeviceModel string 		`json:"deviceModel"`
	VC  string 				`json:"vc"`
	Net string				`json:"net"`
	Ip  string 				`json:"ip"`
	Uid string 				`json:"uid"`
	Sn  string 				`json:"sn"`
	Language string 		`json:"language"`
	Business *int16			`json:"business"`
	RealAppId int64        `json:"rel_app_id"`
}
type PredictReq struct{
	SessionID string 	`json:"session_id"`
	ReqId  string 		`json:"requestId"`
	Imei   string 		`json:"imei"`
	AdLocations []struct{
		PositionId int64 	`jso:"positionId"`
		PageId int64 		`json:"pageId"`
		BlockId int64		`json:"blockId"`
		RankId int64		`json:"rankId"`
		CategoryId int64 	`json:"categoryId"`
		TagId int64 		`json:"tagId"`
		Position int32 		`json:"position"`
		PostionType int32 	`json:"positionType"`
	} `json:"adLocations"`
	DeviceModel string 		`json:"deviceModel"`
	VC  string 				`json:"vc"`
	Net string				`json:"net"`
	Ip  string 				`json:"ip"`
	Uid string 				`json:"uid"`
	Sn  string 				`json:"sn"`
	Language string 		`json:"language"`
	EnableRank *bool 		`json:"enableRank"`
	Topk int16           	`json:"topk"`
	Business *int16			`json:"business"`
}
type Idea_predict struct{
	SessionID string 	`json:"session_id"`
	ReqId  string 		`json:"requestId"`
	Imei   string 		`json:"imei"`
	Ideas []struct{
		Source int16		`json:"source"`
		IdeaId int32		`json:"idea_id"`
	} `json:"ideas"`
	DeviceModel string 		`json:"device"`
	VC  string 				`json:"vc"`
	Net string				`json:"net"`
	Ip  string 				`json:"ip"`
	Uid string 				`json:"uid"`
	Sn  string 				`json:"sn"`
	Language string 		`json:"language"`
}

/*
从上下文中取出callback_req 按照req转换
 */
func getReqFromContext(name string,c *gin.Context) {
	switch name {
	case "callback":
		req := &CallbackReq{}
		if c.BindJSON(req) == nil {
			c.Set("callback_req", req)
		}
		break
	case "predictbinary":
		req := &PredictReq{}
		if c.BindJSON(req) == nil {
			c.Set("predictbinary_req", req)
			fmt.Print("bindjson success\n")
		}
		fmt.Print("this is predictbinary_req struct\n")
		fmt.Print(req.SessionID," ",req.ReqId," ",req.Imei," ",req.Ip," ",req.DeviceModel," ",req.Net," ",req.VC," ")
		fmt.Print(req.AdLocations,"\n")
		break
	case "idea_predict":
		req := &Idea_predict{}
		if c.BindJSON(req) == nil {
			c.Set("idea_predict_req", req)
			//fmt.Print("bindjson success\n")
		}
		//fmt.Print("this is idea_predict_req struct\n")
		//fmt.Print(req.SessionID," ",req.ReqId," ",req.Imei," ",req.Ip," ",req.DeviceModel," ",req.Net," ",req.VC," ")
		//fmt.Print(req.Ideas[1].Source," ",req.Ideas[1].IdeaId,"\n")
		break
	}

	return
	/*
	interface{}接口类型向 普通类型转换需要类型断言  传递一个引用
	*/
}

func initThriftreq(name string, thriftreq *predict.Req, req interface{}) {
	switch name {
	case "callback":
		initCallbackreq(thriftreq,req.(*CallbackReq))
		break
	case "predictbinary":
		initPredictbinary(thriftreq,req.(*PredictReq))
		break
	case "idea_predict":
		initIdea_predict(thriftreq,req.(*Idea_predict))
		break
	}
	return
}


//func get_context_imei_sessionid(c *gin.Context) (string,string,map[string]interface{},[]map[interface{}]interface{},bool,int16,int16){
//	var context_dict = make(map[string]interface{},0)
//	var enable_rank bool
//	var topk int16
//	var business int16
//
//	cacheBody,err := ioutil.ReadAll(c.Request.Body)
//	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(cacheBody))
//	if err == nil {
//		fmt.Print(string(cacheBody))
//	}
//	imei,_ := c.Get("imei")
//	//fmt.Print("\n")
//	//fmt.Print(imei)
//	//fmt.Print("\n")
//	context_dict["net"],_ = c.Get("net")
//	context_dict["device_model"],_ = c.Get("deviceModel")
//	context_dict["ip"],_ = c.Get("ip")
//	context_dict["isp"],_ = c.Get("isp")
//	context_dict["vc"],_ = c.Get("vc")
//	context_dict["language"],_ = c.Get("language")
//	context_dict["uid"],_ = c.Get("uid")
//	context_dict["sn"],_ = c.Get("sn")
//	context_dict["request_id"],_ = c.Get("requestId")
//	session_id,_ := c.Get("sessionId")
//	if session_id == nil{
//		session_id,_ = c.Get("session_id")
//	}else{
//		session_id = ""
//	}
//	adLocations,_ := c.Get("adLocations")
//	adLocations_2 := cast.ToString(adLocations)
//	var adLocations_3 []map[interface{}]interface{}
//	json.Unmarshal([]byte(adLocations_2), &adLocations_3)
//
//
//	if enable_rank_2,exits := c.Get("enableRank");exits{
//		enable_rank = cast.ToBool(enable_rank_2)
//	}
//	enable_rank = true
//	if topk_2,exits := c.Get("topk");exits{
//		topk = cast.ToInt16(topk_2)
//	}
//	topk = 0
//	if business_2,exits:= c.Get("business");exits{
//		business = cast.ToInt16(business_2)
//	}
//	business = 1
//	return cast.ToString(imei), cast.ToString(session_id), context_dict, adLocations_3, enable_rank, topk,business
//}
//
//func get_context_imei_sessionid_2(c *gin.Context) (string,string,map[string]string){
//	var context_dict = make(map[string]string,0)
//	cacheBody,err := ioutil.ReadAll(c.Request.Body)
//	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(cacheBody))
//	var callbackreq Callbackreq
//	json.Unmarshal(cacheBody, &callbackreq)
//	if err == nil {
//		fmt.Print(string(cacheBody))
//	}
//	imei := callbackreq.Imei
//	context_dict["net"] = callbackreq.Net
//	context_dict["device_model"] = callbackreq.DeviceModel
//	context_dict["ip"] = callbackreq.Ip
//	context_dict["isp"] = callbackreq.Isp
//	context_dict["vc"] = callbackreq.Vc
//	context_dict["language"] = callbackreq.Language
//	context_dict["uid"] = callbackreq.Uid
//	context_dict["sn"] = callbackreq.Sn
//	context_dict["request_id"] = callbackreq.RequestId
//	session_id := callbackreq.Sessionid
//	//fmt.Print("\n")
//	//fmt.Print(imei)
//	//fmt.Print("\n")
//	//for key,val := range context_dict{
//	//	fmt.Print(key,"::")
//	//	fmt.Print(val,"\n")
//	//}
//	//fmt.Print("\n")
//	//fmt.Print(session_id)
//	//fmt.Print("\n")
//
//	return cast.ToString(imei), cast.ToString(session_id), context_dict
//}

//func predict_binary(c *gin.Context)   {
//
//	//1。FIXME 待测
//	//_,response := get_check_error_response(c)
//	//if response != nil {
//	//	return
//	//}
//
//	//2。全拿了空   //FIXME
//	imei, session_id, context_dict, adLocations, enable_rank,topk,business := get_context_imei_sessionid(c)
//
//	fmt.Print("\n")
//	fmt.Print(imei, session_id, context_dict, adLocations, enable_rank,topk,business)
//	fmt.Print("\n")
//
//
//	//3  adLocations类型转换问题
//	params_json_string :="{\"3rd_rec\": {\"clk_model_name\": \"meizu_display_beta_0\"}}"
//	is_debug := true
//	is_binary := true
//	//rank_results, version := get_predict_result(business, imei, session_id, adLocations, context_dict, params_json_string, enable_rank, is_debug, is_binary, topk)
//	rank_results, _ := get_predict_result(business, imei, session_id, adLocations, context_dict, params_json_string, enable_rank, is_debug, is_binary, topk)
//
//	//3.
//	//abtestParameters := getABFromReq(c)
//
//	//4. protobuf
//	var final_result string
//	if is_binary {
//		binary_display_result := &toutiao_meizu_app.BinaryDisplayResult{}
//		binary_display_result.Code = SUCCESS
//		*binary_display_result.Message = "success"
//		binary_display_result.Value = rank_results
//		final_result = binary_display_result.String()
//	}
//	//ab_test
//
//	c.JSON(200, final_result)						//返回？
//}

/*需要
1。 c.get 之前上下文中保存好的json数据
2。 对thrift 结构体初始化   为了不传递冗余参数  重新构造对应方法的thrift结构体
3。 CallClient 传递方法名称+已构造的thrift
*/
func callback(c *gin.Context){
	getReqFromContext("callback",c)
	callback_req,_ := c.Get("callback_req")
	thriftreq := predict.NewReq()
	/*
	按需初始化thrift  传递指针
	*/
	initThriftreq("callback",thriftreq,callback_req)
	thriftreq.Debug = true
	/*
	判断thrift 是否被初始化完毕
	*/
	//fmt.Print("初始化相应callback的thrift  \n")
	//fmt.Print(thriftreq.SessionId," ",*thriftreq.ReqId," ",thriftreq.Imei," ",thriftreq.ContextInfo," ",thriftreq.ImprList,"business:",thriftreq.Business)
	//fmt.Print("\n")
	/*
	开始处理URL请求   传递thrift结构
 	*/
	rsp,err := CallClient("callback",thriftreq)
	if err != nil {
		fmt.Print("CallClient error and abort!")
		c.Abort()
		return
	}

	c.Set("rsp", rsp)
}


func predict_binary(c *gin.Context){
	getReqFromContext("predictbinary",c)
	predict_binary_req,_ := c.Get("predictbinary_req")
	thriftreq := predict.NewReq()
	/*
	按需初始化thrift  传递指针
	*/
	initThriftreq("predictbinary",thriftreq,predict_binary_req)
	thriftreq.Debug = true
	/*
	判断thrift 是否被初始化完毕
	*/
	//fmt.Print("初始化相应predictbinary的thrift  \n")
	//fmt.Print(thriftreq.SessionId," ",*thriftreq.ReqId," ",thriftreq.Imei," ",thriftreq.ContextInfo," ","business:",thriftreq.Business)
	//fmt.Print(*thriftreq.AdLocations[0])
	//fmt.Print("\n")
	/*
	开始处理URL请求   传递thrift结构
 	*/
	rsp,err := CallClient("predictbinary",thriftreq)
	if err != nil {
		fmt.Print("CallClient error and abort!")
		c.Abort()
		return
	}

	c.Set("rsp", rsp)
}

func idea_predict (c *gin.Context){
	getReqFromContext("idea_predict",c)
	idea_predict_req,_ := c.Get("idea_predict_req")
	thriftreq := predict.NewReq()
	/*
	按需初始化thrift  传递指针
	*/
	initThriftreq("idea_predict",thriftreq,idea_predict_req)
	thriftreq.Debug = true
	/*
	判断thrift 是否被初始化完毕
	*/
	//fmt.Print("初始化相应idea_predict的thrift  \n")
	//fmt.Print(thriftreq.SessionId," ",*thriftreq.ReqId," ",thriftreq.Imei," ",thriftreq.ContextInfo," ",thriftreq.Ideas)
	//fmt.Print("\n")
	/*
	开始处理URL请求   传递thrift结构
	如果thrift server返回err!=nil abort退出。在responsehandler中返回400
	*/
	rsp,err := CallClient("idea_predict",thriftreq)
	if err != nil {
		fmt.Print("CallClient error and abort!")
		c.Abort()
		return
	}

	c.Set("rsp", rsp)
}