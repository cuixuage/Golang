package main

import (
	"github.com/gin-gonic/gin"
	"meizuapi/thrift_gen/predict"
	"fmt"
	"meizuapi/toutiao_meizu_app"
	"reflect"
	"protobuf/proto"
	"encoding/json"
	"compress/gzip"
	"bytes"
	"io/ioutil"
)

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
		fmt.Print("\n")
		if c.BindJSON(req) == nil {
			c.Set("predictbinary_req", req)
		}
		//fmt.Print("this is predictbinary_req struct\n")
		//fmt.Print(req.SessionID," ",req.ReqId," ",req.Imei," ",req.Ip," ",req.DeviceModel," ",req.Net," ",req.VC," ")
		//fmt.Print(req.Data[0].Sorce,"\n")
		break
	case "recommend":
		req := &RecommendReq{}
		fmt.Print("\n")
		if c.BindJSON(req) == nil {
			c.Set("recommend_req", req)
		}
		break
	case "idea_predict":
		req := &Idea_predict{}
		if c.BindJSON(req) == nil {
			c.Set("idea_predict_req", req)
		}
		break
	case "search":
		req := &Search{}
		fmt.Print("start bind search\n")
		if c.BindJSON(req) == nil {
			c.Set("search_req", req)
			fmt.Print("bindjson success\n")
		}
		break
	}
	return
}

func callback(c *gin.Context){
	getReqFromContext("callback",c)
	callback_req,_ := c.Get("callback_req")
	thriftreq := predict.NewReq()
	initCallbackreq(thriftreq,callback_req.(*CallbackReq))
	debug,exits := c.Get("debug")
	if exits && debug == true{
		thriftreq.Debug = true
	}
	_,err := CallClient("callback",thriftreq)
	if err != nil {
		fmt.Print("CallClient error and abort!")
		c.Abort()
		return
	}
	c.Set("pb", "false")
}


func predict_binary(c *gin.Context){
	getReqFromContext("predictbinary",c)
	predict_binary_req,_ := c.Get("predictbinary_req")
	thriftreq := predict.NewReq()
	initPredictbinary(thriftreq,predict_binary_req.(*PredictReq))
	debug,exits := c.Get("debug")
	if exits && debug == true{
		thriftreq.Debug = true
	}
	rsp,err := CallClient("predictbinary",thriftreq)
	if err != nil {
		fmt.Print("\n")
		fmt.Print(err,"\n")
		fmt.Print("CallClient error and abort!\n")
		c.Abort()
		return
	}
	pb_result := &toutiao_meizu_app.BinaryDisplayResult{}
	pb_result.Code = SUCCESS
	pb_result.Message = &StrSuccess
	r := reflect.ValueOf(rsp)
	f_rankresult := reflect.Indirect(r).FieldByName("RankResults")
	pb_result.Value = f_rankresult.Bytes()
	pb_bytes,_ := proto.Marshal(pb_result)
	c.Data(200,"application/octet-stream",pb_bytes)
	c.Set("rsp", rsp)
}

func recommend_predict(c *gin.Context){
	getReqFromContext("recommend",c)
	Recommend_req,_ := c.Get("recommend_req")
	thriftreq := predict.NewReq()
	initRecommend(thriftreq,Recommend_req.(*RecommendReq))
	debug,exits := c.Get("debug")
	if exits && debug == true {
		thriftreq.Debug = true
	}
	rsp,err := CallClient("recommend",thriftreq)
	if err != nil {
		//fmt.Print("CallClient error and abort!\n")
		c.Abort()
		return
	}
	c.Set("pb", "true")
	c.Set("rsp",rsp)
}



func idea_predict (c *gin.Context){
	getReqFromContext("idea_predict",c)
	idea_predict_req,_ := c.Get("idea_predict_req")
	thriftreq := predict.NewReq()
	initIdea_predict(thriftreq,idea_predict_req.(*Idea_predict))
	debug,exits := c.Get("debug")
	if exits && debug == true{
		thriftreq.Debug = true
	}
	rsp,err := CallClient("idea_predict",thriftreq)
	if err != nil {
		fmt.Print("CallClient error and abort!")
		c.Abort()
		return
	}
	c.Set("rsp", rsp)
	c.Set("pb","true")
}

func search_pb (c *gin.Context){
	getReqFromContext("search",c)
	idea_predict_req,_ := c.Get("search_req")
	thriftreq := predict.NewReq()
	initSearch(thriftreq,idea_predict_req.(*Search))
	debug,exits := c.Get("debug")
	if exits && debug == true{
		thriftreq.Debug = true
	}
	rsp,err := CallClient("search",thriftreq)
	if err != nil {
		fmt.Print("CallClient error and abort!")
		c.Abort()
		return
	}
	c.Set("rsp", rsp)
	c.Set("pb","true")
}


func upload_behavior (c *gin.Context){
	thriftreq := predict.NewAckReq()
	upMetaData := &Upload_Meta{}
	json.Unmarshal([]byte(c.Request.Header.Get("METADATA")), &upMetaData)
	thriftreq.MessageType = &upMetaData.Dataname
	debug,exits := c.Get("debug")
	if exits && debug == true{
		fmt.Print("\n")
		fmt.Print("Loger: upload_behavior_dataname is  ",upMetaData.Dataname,"\n")
	}
	body,exits := c.Get("cachebody")
	if exits {
		strbody, err := body.([]byte)
		if err == false {
			c.Abort()
			return
		}
		buf := bytes.NewBuffer(strbody)
		req, _ := gzip.NewReader(buf)
		if req != nil{
			byteReq, _ := ioutil.ReadAll(req) //打开
			ioutil.NopCloser(req)
			req.Close()
			thriftreq.Message = string(byteReq)
		}
	}
	_,err2 := CallClient("upload_behavior",thriftreq)
	if err2 != nil {
		fmt.Print("CallClient error and abort!")
		c.Abort()
		return
	}
	c.Set("pb","upload")
	c.Set("uploadReqId",upMetaData.ReqId)
}


func upload_applist (c *gin.Context){
	thriftreq := predict.NewAckReq()
	upMetaData := &Upload_Meta{}
	json.Unmarshal([]byte(c.Request.Header.Get("METADATA")), &upMetaData)
	debug,exits := c.Get("debug")
	if exits && debug == true{
		fmt.Print("\n")
		fmt.Print("Loger: upload_applist_dataname is  ",upMetaData.Dataname,"\n")
	}
	body,exits := c.Get("cachebody")
	if exits {
		strbody, err := body.([]byte)
		if err == false {
			c.Abort()
			return
		}
		buf := bytes.NewBuffer(strbody)
		req, _ := gzip.NewReader(buf)
		if req != nil{
			byteReq, _ := ioutil.ReadAll(req)
			ioutil.NopCloser(req)
			req.Close()
			thriftreq.Message = string(byteReq)
		}
	}
	_,err2 := CallClient("upload_applist",thriftreq)
	if err2 != nil {
		fmt.Print("CallClient error and abort!")
		c.Abort()
		return
	}
	c.Set("pb","upload")
	c.Set("uploadReqId",upMetaData.ReqId)
}

func idea_predict_test(c *gin.Context){
	c.Set("debug",true)
	idea_predict(c)
}
func recommend_predict_test(c *gin.Context){
	c.Set("debug",true)
	recommend_predict(c)
}
func callback_test(c *gin.Context){
	c.Set("debug",true)
	callback(c)
}
func predict_binary_test(c *gin.Context){
	c.Set("debug",true)
	predict_binary(c)
}
func search_pb_test(c *gin.Context){
	c.Set("debug",true)
	search_pb(c)
}
func upload_behavior_test(c *gin.Context){
	c.Set("debug",true)
	upload_behavior(c)
}
func upload_applist_test(c *gin.Context){
	c.Set("debug",true)
	upload_applist(c)
}
