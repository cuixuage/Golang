package main

import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
)

type Response struct{
	Status string `json:"status"`
	RankResult byte `json:"rankresults"`
}
//callback
func httpPostForm_callback() {
	url := "http://127.0.0.1:8000/collaborate/meizu/predict/callback"
	fmt.Println("URL:", url)

	post :=`{"session_id":"","requestId":"a3004ff101610a03053a27100093f7d9","imei":"f3e6f5d3637e708ed544dee4be66d289f6f706d7a47dc8d3423cd42d737847b7","data":[{"score":0.032694172114133835,"idea_id":111771}],"deviceModel":"M5710","vc":"","net":"0","ip":"182.39.210.31"}`;
	fmt.Println( "post", post)
	var jsonStr = []byte(post)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP_METADATA", `{"checkInfo":"45950f6676db04492b2ac431ca596d75","nonce":1518071142348,"appkey":"c7d805ab7f396ed6d420cd8da0f788ca","reqid":"542074170161ac10b12b27100000004e"}`)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	//defer resp.Body.Close()
	//
	fmt.Println("Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Body:", string(body),"\n")


	//r := reflect.ValueOf(body)
	//f_status := reflect.Indirect(r).FieldByName("Status")
	//f_rankresult := reflect.Indirect(r).FieldByName("RankResults")
	//status := f_status.String()
	//fmt.Print("checkresponse status is:",status,"\n")
	//fmt.Print(f_rankresult.Bytes())

}

//idea_predict
func httpPostForm_idea_predict() {
	url := "http://127.0.0.1:8000/collaborate/meizu/idea_predict"
	fmt.Println("URL:", url)

	post :=`{"ideas":[{"source":0,"idea_id":72948},{"source":0,"idea_id":67394},{"source":0,"idea_id":72982},{"source":0,"idea_id":72981},{"source":0,"idea_id":72971},{"source":0,"idea_id":67572},{"source":0,"idea_id":73033},{"source":0,"idea_id":73028},{"source":0,"idea_id":73031},{"source":0,"idea_id":73012},{"source":0,"idea_id":73073},{"source":0,"idea_id":73071},{"source":0,"idea_id":73061},{"source":0,"idea_id":73054},{"source":0,"idea_id":73067},{"source":0,"idea_id":73068},{"source":0,"idea_id":73043},{"source":0,"idea_id":73050},{"source":0,"idea_id":73051},{"source":0,"idea_id":72851},{"source":0,"idea_id":72829},{"source":0,"idea_id":72830},{"source":0,"idea_id":72833},{"source":0,"idea_id":72826},{"source":0,"idea_id":72827},{"source":0,"idea_id":72818},{"source":0,"idea_id":72817},{"source":0,"idea_id":72806},{"source":0,"idea_id":72805},{"source":0,"idea_id":72897},{"source":0,"idea_id":72898},{"source":0,"idea_id":72894},{"source":0,"idea_id":72882},{"source":0,"idea_id":67749},{"source":0,"idea_id":72868},{"source":0,"idea_id":72875},{"source":0,"idea_id":70382},{"source":0,"idea_id":72856},{"source":0,"idea_id":72854},{"source":0,"idea_id":72857},{"source":0,"idea_id":72928},{"source":0,"idea_id":72913},{"source":0,"idea_id":72642},{"source":0,"idea_id":72641}],"imei":"4ff53e8501a343c0b5b9f6269f5c27833f13154fba1e01bdbfd4e66482ef5b1b","device":"M8810","vc":"","net":"0","ip":"58.30.17.16"}`;
     	fmt.Println( "post", post)
	var jsonStr = []byte(post)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP_METADATA", `{"checkInfo":"a906f44f3ea8933da71143244c2e944b","nonce":1518074588726,"appkey":"c7d805ab7f396ed6d420cd8da0f788ca","reqid":"ee2a744b0161ac1131ab271000017b35"}`)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	//defer resp.Body.Close()
	//
	fmt.Println("Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Body:", string(body),"\n")


	//r := reflect.ValueOf(body)
	//f_status := reflect.Indirect(r).FieldByName("Status")
	//f_rankresult := reflect.Indirect(r).FieldByName("RankResults")
	//status := f_status.String()
	//fmt.Print("checkresponse status is:",status,"\n")
	//fmt.Print(f_rankresult.Bytes())

}


//predict_binary
func httpPostForm_predict_binary() {
	url := "http://127.0.0.1:8000/collaborate/meizu/predict_binary"
	fmt.Println("URL:", url)

	post :=`{"imei":"00007792e83bea48158bf45ac41e2540e3daf14ad6e351729f4262b8eaaede6f","session_id":"", "adLocations":[{"positionId":26,"pageId":5000,"blockId": 10000,"rankId": 9920,"categoryId": 0,"tagId": 0,"position": 3,"positionType": 1}], "enableRank":true, "topk":10, "deviceModel":"M5710","vc":"","net":"0","ip":"182.39.210.31"}`;
	fmt.Println( "post", post)
	var jsonStr = []byte(post)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP_METADATA", `{"checkInfo":"a906f44f3ea8933da71143244c2e944b","nonce":1518074588726,"appkey":"c7d805ab7f396ed6d420cd8da0f788ca","reqid":"ee2a744b0161ac1131ab271000017b35"}`)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	//defer resp.Body.Close()
	//
	fmt.Println("Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Body:", string(body),"\n")


	//r := reflect.ValueOf(body)
	//f_status := reflect.Indirect(r).FieldByName("Status")
	//f_rankresult := reflect.Indirect(r).FieldByName("RankResults")
	//status := f_status.String()
	//fmt.Print("checkresponse status is:",status,"\n")
	//fmt.Print(f_rankresult.Bytes())

}

func main(){
	httpPostForm_callback()
	//httpPostForm_idea_predict()
	//httpPostForm_predict_binary()
}
