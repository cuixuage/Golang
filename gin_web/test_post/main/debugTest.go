package main

import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
	"meizuapi/toutiao_meizu_app"
	"protobuf/proto"
)

//callback
func debug_callback() {
	//url := "http://10.8.64.231:4640/collaborate/meizu/predict/callback_test/"
	url := "http://127.0.0.1:8000/collaborate/meizu/predict/callback_test/"
	fmt.Println("URL:", url)

	post :=`{"session_id":"","requestId":"a3004ff101610a03053a27100093f7d9","imei":"f3e6f5d3637e708ed544dee4be66d289f6f706d7a47dc8d3423cd42d737847b7","data":[{"score":0.032694172114133835,"idea_id":111771}],"deviceModel":"M5710","vc":"","net":"0","ip":"182.39.210.31"}`;
	//fmt.Println( "post", post)
	var jsonStr = []byte(post)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("METADATA", `{"appkey":"c7d805ab7f396ed6d420cd8da0f788ca","bizModule":"meizu_uxip","checkInfo":"71fad6afeb17abc7a20f6001cda0e3f5","contentFormat":"batch","data_count":1000,"dataname":"meizu_userdata_thrid_cp_export_activated","nonce":1518163771154,"reqid":"bf01799c01610a038c6b2710066045cb"}`)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	//defer resp.Body.Close()
	//
	fmt.Println("Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Print("body is ",string(body),"\n")
}

func debug_test_idea_predict(){
	url := "http://10.8.64.231:4640/collaborate/meizu/idea_predict_test/"
	url_2 := "http://127.0.0.1:8000/collaborate/meizu/idea_predict_test/"
	post :=`{"ideas":[{"source":0,"idea_id":72948},{"source":0,"idea_id":67394},{"source":0,"idea_id":72982},{"source":0,"idea_id":72981},{"source":0,"idea_id":72971},{"source":0,"idea_id":67572},{"source":0,"idea_id":73033},{"source":0,"idea_id":73028},{"source":0,"idea_id":73031},{"source":0,"idea_id":73012},{"source":0,"idea_id":73073},{"source":0,"idea_id":73071},{"source":0,"idea_id":73061},{"source":0,"idea_id":73054},{"source":0,"idea_id":73067},{"source":0,"idea_id":73068},{"source":0,"idea_id":73043},{"source":0,"idea_id":73050},{"source":0,"idea_id":73051},{"source":0,"idea_id":72851},{"source":0,"idea_id":72829},{"source":0,"idea_id":72830},{"source":0,"idea_id":72833},{"source":0,"idea_id":72826},{"source":0,"idea_id":72827},{"source":0,"idea_id":72818},{"source":0,"idea_id":72817},{"source":0,"idea_id":72806},{"source":0,"idea_id":72805},{"source":0,"idea_id":72897},{"source":0,"idea_id":72898},{"source":0,"idea_id":72894},{"source":0,"idea_id":72882},{"source":0,"idea_id":67749},{"source":0,"idea_id":72868},{"source":0,"idea_id":72875},{"source":0,"idea_id":70382},{"source":0,"idea_id":72856},{"source":0,"idea_id":72854},{"source":0,"idea_id":72857},{"source":0,"idea_id":72928},{"source":0,"idea_id":72913},{"source":0,"idea_id":72642},{"source":0,"idea_id":72641}],"imei":"4ff53e8501a343c0b5b9f6269f5c27833f13154fba1e01bdbfd4e66482ef5b1b","device":"M8810","vc":"","net":"0","ip":"58.30.17.16"}`;
	fmt.Println( "post", post)
	post2 :=`{"ideas":[{"source":0,"idea_id":72948},{"source":0,"idea_id":67394},{"source":0,"idea_id":72982},{"source":0,"idea_id":72981},{"source":0,"idea_id":72971},{"source":0,"idea_id":67572},{"source":0,"idea_id":73033},{"source":0,"idea_id":73028},{"source":0,"idea_id":73031},{"source":0,"idea_id":73012},{"source":0,"idea_id":73073},{"source":0,"idea_id":73071},{"source":0,"idea_id":73061},{"source":0,"idea_id":73054},{"source":0,"idea_id":73067},{"source":0,"idea_id":73068},{"source":0,"idea_id":73043},{"source":0,"idea_id":73050},{"source":0,"idea_id":73051},{"source":0,"idea_id":72851},{"source":0,"idea_id":72829},{"source":0,"idea_id":72830},{"source":0,"idea_id":72833},{"source":0,"idea_id":72826},{"source":0,"idea_id":72827},{"source":0,"idea_id":72818},{"source":0,"idea_id":72817},{"source":0,"idea_id":72806},{"source":0,"idea_id":72805},{"source":0,"idea_id":72897},{"source":0,"idea_id":72898},{"source":0,"idea_id":72894},{"source":0,"idea_id":72882},{"source":0,"idea_id":67749},{"source":0,"idea_id":72868},{"source":0,"idea_id":72875},{"source":0,"idea_id":70382},{"source":0,"idea_id":72856},{"source":0,"idea_id":72854},{"source":0,"idea_id":72857},{"source":0,"idea_id":72928},{"source":0,"idea_id":72913},{"source":0,"idea_id":72642},{"source":0,"idea_id":72641}],"imei":"4ff53e8501a343c0b5b9f6269f5c27833f13154fba1e01bdbfd4e66482ef5b1b","device":"M8810","vc":"","net":"0","ip":"58.30.17.20"}`;
	fmt.Println( "post", post)
	fmt.Println("URL:", url)
	res1 := debug_idea_predict(url,post)
	res2 := debug_idea_predict(url_2,post2)
	fmt.Print("\n")
	fmt.Print("\n")
	fmt.Print("\n")
	for key,_ := range res1{
		//fmt.Print(key," ",*res1[key].IdeaId-*res2[key].IdeaId," ",*res1[key].Ctr-*res2[key].Ctr," ",*res1[key].Cvr-*res2[key].Cvr)
		fmt.Print(key," ",*res2[key].IdeaId," ",*res1[key].Ctr," ",*res2[key].Ctr," ",*res1[key].Ctr-*res2[key].Ctr," \n")
	}

}
//idea_predict
func debug_idea_predict(url string,post string) [] *toutiao_meizu_app.IdeaCTR {
	var jsonStr = []byte(post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("METADATA", `{"checkInfo":"a906f44f3ea8933da71143244c2e944b","nonce":1518074588726,"appkey":"c7d805ab7f396ed6d420cd8da0f788ca","reqid":"ee2a744b0161ac1131ab271000017b35"}`)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	//fmt.Println("Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)

	rsp := &toutiao_meizu_app.BinaryDisplayResult{}
	proto.Unmarshal(body,rsp)
	fmt.Print("",*rsp.Code," ",*rsp.Message," ",rsp.Value,"\n")
	rsp_2_rankresult := &toutiao_meizu_app.IdeaCTRList{}
	proto.Unmarshal(rsp.Value,rsp_2_rankresult)
	fmt.Print("rankresult version: ",*rsp_2_rankresult.Version,"\n")
	fmt.Print("rankresult len:",len(rsp_2_rankresult.Data))
	res := rsp_2_rankresult.Data
	//for key,val := range rsp_2_rankresult.Data{
	//	fmt.Print(key," ",*val.IdeaId," ",*val.Ctr," ",val.Cvr,"\n")
	//}
	return res
}


func debug_binary_predict() {
	url := "http://10.8.64.231:4640/collaborate/meizu/predict_binary_test/"
	//url_2 := "http://127.0.0.1:8000/collaborate/meizu/predict_binary_test/"
	post := `{"devicemodel":"M6820","source":1,"sessionId":"","imei":"343d6e0796e7bf2bb3efc14149520a131690a135cd5b083075c41c0617009c6b","vc":6019010,"topk":0,"enableRank":false,"net":"wifi","adLocations":[{"position":1,"appId":0,"blockId":12172,"positionId":1950,"pageId":6028,"categoryId":0,"positionType":1,"rankId":7000},{"position":2,"appId":0,"blockId":12172,"positionId":1951,"pageId":6028,"categoryId":0,"positionType":1,"rankId":7000},{"position":3,"appId":0,"blockId":12172,"positionId":2021,"pageId":6028,"categoryId":0,"positionType":1,"rankId":7000}],"ip":"49.90.155.191"}`
	//res1 := debug_predict_binary(url, post)
	debug_predict_binary(url, post)
	//res2 := debug_predict_binary(url_2, post)
	fmt.Print("\n")
	fmt.Print("\n")
	fmt.Print("\n")
	//for key, _ := range res1 {
	//	fmt.Print(key, " ", *res1[key].PositionId - *res2[key].PositionId, " ", *res1[key].Ctr - *res2[key].Ctr, " ",
	//		*res1[key].UnitId - *res2[key].UnitId, " ", *res1[key].AppId - *res2[key].AppId, " \n")
	//}
	//for key, _ := range res1 {
	//	fmt.Print(key, " ", *res1[key].PositionId , " ", *res1[key].Ctr ," ",
	//		*res1[key].UnitId, " ", *res1[key].AppId , " \n")
	//}
	//fmt.Print(res2)
}

//predict_binary
func debug_predict_binary(url string,post string) []*toutiao_meizu_app.DisplayUnit{
	fmt.Println("URL:", url)
	//fmt.Println( "post", post)

	var jsonStr = []byte(post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("METADATA", `{"checkInfo":"a906f44f3ea8933da71143244c2e944b","nonce":1518074588726,"appkey":"c7d805ab7f396ed6d420cd8da0f788ca","reqid":"ee2a744b0161ac1131ab271000017b35"}`)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	//defer resp.Body.Close()
	//
	fmt.Println("Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(body))
	rsp := &toutiao_meizu_app.BinaryDisplayResult{}
	proto.Unmarshal(body,rsp)
	fmt.Print("",*rsp.Code," ",*rsp.Message," ",rsp.Value,"\n")

	rsp_2_rankresult := &toutiao_meizu_app.DisplayUnitList{}
	proto.Unmarshal(rsp.Value,rsp_2_rankresult)
	fmt.Print("rankresult version: ",*rsp_2_rankresult.Version,"\n")
	fmt.Print("rankresult len:",len(rsp_2_rankresult.Data),"\n")
	res := rsp_2_rankresult.Data
	return res

}


func debug_test_search(){
	url := "http://10.8.64.231:4640/collaborate/meizu/search_test/"
	url2 := "http://127.0.0.1:8000/collaborate/meizu/search_pb_test/"
	post :=`{"uid":150436900,"kw":"dnf","imei":"ab722b5dcefff53c733c9ad5ca94b655be70729a38e5c9e253139fef9872046c","net":"wifi",
    "ip":"183.28.28.92",
    "devicemodel":"M6887",
    "source":1,
    "sessionId":"959432243dd9fea89ae5ddb1c1bdf2e4",
    "vc":6019010,
    "enableRank":false,
    "language":"zh-CN"
}`
	res1 := debug_search(url,post)
	res2 :=debug_search(url2,post)
	for key,_ := range res1{
		fmt.Print(key," ",*res1[key].AppId-*res2[key].AppId," ",*res1[key].UnitId-*res1[key].UnitId," ",
			*res1[key].Ctr-*res2[key].Ctr," ",*res1[key].Relevance-*res2[key].Relevance,"\n")
	}
}

func debug_search(url string,post string) []*toutiao_meizu_app.SearchUnit{
	fmt.Println("URL:", url)
	//fmt.Println( "post", post)
	var jsonStr = []byte(post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("METADATA", `{"checkInfo":"a906f44f3ea8933da71143244c2e944b","nonce":1518074588726,"appkey":"c7d805ab7f396ed6d420cd8da0f788ca","reqid":"ee2a744b0161ac1131ab271000017b35"}`)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	//defer resp.Body.Close()
	//
	fmt.Println("Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Print(string(body))
	rsp := &toutiao_meizu_app.BinaryDisplayResult{}
	proto.Unmarshal(body,rsp)
	fmt.Print("",*rsp.Code," ",*rsp.Message," ",rsp.Value,"\n")
	//
	rsp_2_rankresult := &toutiao_meizu_app.SearchUnitList{}
	proto.Unmarshal(rsp.Value,rsp_2_rankresult)
	fmt.Print("rankresult version: ",*rsp_2_rankresult.Version,"\n")
	fmt.Print("rankresult len:",len(rsp_2_rankresult.Data),"\n")
	res := rsp_2_rankresult.Data
	return res

}

func debug_Upload_behavior() {
	//url := "http://10.8.64.231:4640/collaborate/meizu/upload_behavior_test/"
	url := "http://127.0.0.1:8000/collaborate/meizu/upload_behavior_test/"
	post :=""
	fmt.Println("URL:", url)
	//fmt.Println( "post", post)
	var jsonStr = []byte(post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("METADATA", `{"appkey":"c7d805ab7f396ed6d420cd8da0f788ca","bizModule":"meizu_uxip","checkInfo":"7679c44f245b2fe9eb1f0182f1222333","contentFormat":"batch","data_count":1000,"dataname":"meizu_userdata_thrid_cp_export_exposure","nonce":"1518163772262","reqid":"8a87064c-9581-487c-a5e3-55ca5c8a180e"}`)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(body))
}

func debug_Upload_applist() {
	//url := "http://10.8.64.231:4640/collaborate/meizu/upload_applist_test/"
	url := "http://127.0.0.1:8000/collaborate/meizu/upload_applist_test/"
	post :=""
	fmt.Println("URL:", url)
	//fmt.Println( "post", post)
	var jsonStr = []byte(post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("METADATA", `{"appkey":"c7d805ab7f396ed6d420cd8da0f788ca","bizModule":"meizu_uxip","checkInfo":"ed08e6c247a88723aa25bc7cf2852235","contentFormat":"batch","data_count":1000,"dataname":"meizu_appcenter_user_behavior_export_install","nonce":"1518163775658","reqid":"4ed79e56-8079-4fae-99e8-27af115d5afe"}`)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(body))
}



func debug_Test_Recommend_predict() {
	url := "http://10.8.64.231:4640/collaborate/meizu/recommend/normal_test"
	url_2 := "http://127.0.0.1:8000/collaborate/meizu/recommend/normal_test"
	post := `{"appId":1942736,"imei":"0621260225f8af8b58c5c9da6804a44dd8260d8583ca153eae3b5d516ea66d77","topk":0,"adLocations":[{"position":1,"appId":0,"blockId":11782,"positionId":1805,"pageId":5007,"categoryId":0,"positionType":1,"rankId":5333},{"position":2,"appId":0,"blockId":11782,"positionId":1804,"pageId":5007,"categoryId":0,"positionType":1,"rankId":5333}],"net":"wifi","ip":"112.5.248.170","devicemodel":"M6210","source":1,"sessionId":"dd8507febe20e4c82079e36d558defc8","vc":6019010,"enableRank":false}`
	res1 := debug_Recommend_predict(url, post)
	res2 := debug_Recommend_predict(url_2, post)
	fmt.Print("\n")
	fmt.Print("\n")
	fmt.Print("\n")
	for key, _ := range res1 {
		fmt.Print(key, " ", *res1[key].PositionId , " ", *res1[key].Ctr - *res2[key].Ctr, " ",
			*res1[key].UnitId - *res2[key].UnitId, " ", *res1[key].AppId - *res2[key].AppId, " \n")
	}
}

//Recommend_predict
func debug_Recommend_predict(url string,post string) []*toutiao_meizu_app.DisplayUnit{
	fmt.Println("URL:", url)
	//fmt.Println( "post", post)

	var jsonStr = []byte(post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("METADATA", `{"checkInfo":"22fc4551191d4732f4895f9ecd3394cc","nonce":1518163770397,"appkey":"c7d805ab7f396ed6d420cd8da0f788ca","reqid":"a696f42172c71e35206e72880954c2aa"}`)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	rsp := &toutiao_meizu_app.BinaryDisplayResult{}
	proto.Unmarshal(body,rsp)
	fmt.Print("",*rsp.Code," ",*rsp.Message," ",rsp.Value,"\n")
	rsp_2_rankresult := &toutiao_meizu_app.DisplayUnitList{}
	proto.Unmarshal(rsp.Value,rsp_2_rankresult)
	fmt.Print("rankresult version: ",*rsp_2_rankresult.Version,"\n")
	fmt.Print("rankresult len:",len(rsp_2_rankresult.Data),"\n")
	res := rsp_2_rankresult.Data
	return res

}





func main(){
	//httpPostForm_callback()
	//httpPostForm_idea_predict()
	//httpPostForm_predict_binary()



	//debug_Test_Recommend_predict()					//ok
	//debug_test_search()								//FIXME  search_pb_test 只返回一个结果  search_test路径即可
	//debug_test_idea_predict()							//test第一遍有差距 后面全部一样 是缓存的原因吗？？
	//debug_binary_predict()					//error  && predict_binary_test网关错误
	debug_callback()					//ok
	debug_Upload_behavior()       				//ok??   如何测试post不为空的情况
	debug_Upload_applist()						//没有post 数据

}


// predict_binary error
// idea_predict 有差距
// search_pb commend callback ok
// upload 如何测试post数据
// upload 如何测试post数据