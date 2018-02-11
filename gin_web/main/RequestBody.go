package main

/*
格式标注还挺严格的  json
后面判断是否为空的参数  转换为指针再判断nil
*/
/*
interface{}接口类型向 普通类型转换需要类型断言  传递一个引用
*/
/*
//FIXME upload_behavoir nonce是string
 */
type CheckRequst struct{
	ReqId string `json:"reqid"`
	Appkey string 	`json:"appkey"`
	Checkinfo string `json:"checkInfo"`
	Nonce interface{} `json:"nonce"`				//int or string
}
type JsonResponse struct {
	Code int `json:"code"`
	Value interface{} `json:"value"`
	Message string `json:"message"`
	Redirect string `json:"redirect"`
}
type JsonResponseUpload struct {
	Code int `json:"code"`
	Message string `json:"msg"`
	ReqId string `json:"reqid"`
}
type CallbackReq struct{
	SessionID string 	`json:"session_id"`
	ReqId  string 		`json:"requestId"`
	Imei   string 		`json:"imei"`
	Query  string 		`json:"query"`
	Data []struct{
		PositionId int64 	`jso:"position_id"`
		PageId int64 		`json:"page_id"`
		BlockId int64		`json:"blockId"`
		RankId int64		`json:"rankId"`
		CategoryId int64 	`json:"categoryId"`
		TagId int64 		`json:"tagId"`
		AppId int64			`json:"app_id"`
		UnitId int64		`json:"unit_id"`
		Score float64		`json:"score"`
		IdeaId int32		`json:"idea_id"`
		Keyword string 		`json:"kw"`
		Position int32 		`json:"position"`
		PostionType int32 	`json:"positionType"`
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
	DeviceModel string 	`json:"devicemodel"`
	SessionID string 	`json:"sessionId"`
	ReqId  string 		`json:"requestId"`
	Imei   string 		`json:"imei"`
	AdLocations []struct{
		Position int32 		`json:"position"`
		AppId  int64 		`json:"appId"`
		BlockId int64		`json:"blockId"`
		PositionId int64 	`jso:"positionId"`
		PageId int64 		`json:"pageId"`
		CategoryId int64 	`json:"categoryId"`
		PostionType int32 	`json:"positionType"`
		RankId int64		`json:"rankId"`
		TagId int64 		`json:"tagId"`
	} `json:"adLocations"`
	VC  int64 				`json:"vc"`
	Net string				`json:"net"`
	Ip  string 				`json:"ip"`
	Uid string 				`json:"uid"`
	Sn  string 				`json:"sn"`
	Language string 		`json:"language"`
	Source int16			`json:"source"`
	EnableRank *bool 		`json:"enableRank"`
	Topk *int16           	`json:"topk"`
	Business *int16			`json:"business"`
}
type RecommendReq struct{
	DeviceModel string 	`json:"devicemodel"`
	SessionID string 	`json:"sessionId"`
	ReqId  string 		`json:"requestId"`
	Imei   string 		`json:"imei"`
	AdLocations []struct{
		Position int32 		`json:"position"`
		//AppId  int64 		`json:"appId"`
		BlockId int64		`json:"blockId"`
		PositionId int64 	`jso:"positionId"`
		PageId int64 		`json:"pageId"`
		CategoryId int64 	`json:"categoryId"`
		PostionType int32 	`json:"positionType"`
		RankId int64		`json:"rankId"`
		TagId int64 		`json:"tagId"`
	} `json:"adLocations"`
	VC  int64 				`json:"vc"`
	Net string				`json:"net"`
	Ip  string 				`json:"ip"`
	Uid string 				`json:"uid"`
	Sn  string 				`json:"sn"`
	Language string 		`json:"language"`
	Source int16			`json:"source"`
	EnableRank *bool 		`json:"enableRank"`
	Topk *int16           	`json:"topk"`
	Business *int16			`json:"business"`
	AppId *int64				`json:"appId"`
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
type Search struct{
	SessionID string 	`json:"sessionId"`
	ReqId  string 		`json:"requestId"`
	Imei   string 		`json:"imei"`
	DeviceModel string 		`json:"devicemodel"`
	VC  int32 				`json:"vc"`
	Net string				`json:"net"`
	Ip  string 				`json:"ip"`
	Uid int64 				`json:"uid"`
	Sn  string 				`json:"sn"`
	Language string 		`json:"language"`
	EnableRank *bool 		`json:"enableRank"`
	Source *int16           `json:"source"`
	Business *int16			`json:"business"`
	Kw *string 				`json:"kw"`
}

type Upload_Meta struct{
	Dataname string `json:"dataname"`
	ReqId string 	`json:"reqid"`
}
