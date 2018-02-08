package main

import (
	"runtime"
	"fmt"
	"code.byted.org/gin/ginex"
	"net/http"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"code.byted.org/gin/ginex/ctx"
)

/*
//FIXME requstID定的格式也不一样
 */
type CheckRequst struct{
	ReqId string `json:"reqid"`
	Appkey string 	`json:"appkey"`
	Checkinfo string `json:"checkInfo"`
	Nonce int64 `json:"nonce"`
}
type MeizuzResponse struct {
	Code int `json:"code"`
	Value interface{} `json:"value"`
	Message string `json:"message"`
}
type recoverWriter struct{}
func (rw *recoverWriter) Write(p []byte) (int, error) {
	return gin.DefaultErrorWriter.Write(p)
}
/*
1。 json 序列化操作  将数据按照json格式解析到req中
2。 c.set 在gin上下文中定义变量 vivo_req:req
3。 c.Next() 处理url请求
 */
func JsonRequestContext() gin.HandlerFunc {
		return func(c *gin.Context) {
			if  c.Request.Method != "GET" {
				cacheBody, err := ioutil.ReadAll(c.Request.Body)              //打开
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(cacheBody)) //无操作关闭
					if err == nil {
						c.Set("cachebody", cacheBody)
					}
				fmt.Print("\n")
				fmt.Print("this is cachebody:",string(cacheBody))
				fmt.Print("\n")

				checkreq := &CheckRequst{}
				err = json.Unmarshal([]byte(c.Request.Header.Get("HTTP_METADATA")), &checkreq)
				/*
				header 校验。失败则abort
				 */
				_,e := get_check_error_response(checkreq)
				if e != nil{
					c.Abort()
					c.Next()
				}else{
					fmt.Print("header 校验成功 \n")
				}
				c.Next() //处理url请求  Next should be used only inside middleware
		}
	}
}

/*
1。 c.Next()	 处理url请求
关键：  此请求意外终止后 判断imei是否被锁定（锁定则意味着被其他routine处理）
2。 c.Get("vivo_rsp")  之前在调用url方法中 c。set(vivo_rsp)
3。 c.json 返回
 */
func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		var res *MeizuzResponse = &MeizuzResponse{}
		/*
		注意处理abort情况
		 */
		if c.IsAborted(){
			res.Code = 500
			res.Message = "catch error"
			c.AbortWithStatusJSON(200,res)
		}else if c.Request.Method != "GET" {
			res.Code = 1000
			val, _ := c.Get("rsp")				//就是thrift返回的数据
			res.Value = val
			//fmt.Print("this is callback_rsp = ",val)
			//fmt.Print(res.Value)
			c.JSON(200, val)
		}
	}
}


func main(){
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	//r := gin.Default()

	ginex.Init()
	r := ginex.New()
	InitThriftClient()
	r.Use(gin.RecoveryWithWriter(&recoverWriter{}))
	r.Use(ctx.Ctx())
	r.Use(JsonRequestContext())
	r.Use(ResponseHandler())
	//r.Use(accesslog.AccessLog(accessLogger))
	//if appConf.EnableMetrics && ginex.PSM() != "" {
	//	r.Use(Metrics(ginex.PSM()))
	//}
	//r.Use(BlockImei())


	meizu_router := r.Group("/collaborate/meizu")
	{
		meizu_router.POST("/idea_predict", idea_predict)
		//meizu_router.POST("/recommend/normal", recommend_predict)
		//
		meizu_router.POST("/predict/callback", callback)
		//meizu_router.POST("/predict/callback_test", callback_test)
		//meizu_router.POST("/predict/boost", boost)

		meizu_router.POST("/predict_binary", predict_binary)
		//meizu_router.POST("/predict_binary_test", predict_binary_test)

		//meizu_router.POST("/predict_pb", predict_pb)
		//meizu_router.POST("/search_pb", search_pb)
		//meizu_router.POST("/predict_pb_test", predict_pb_test)
		//meizu_router.POST("/search_pb_test", search_pb_test)
	}

	//meizu_router_upload := r.Group("/collaborate/meizu")
	//{
	//	meizu_router_upload.POST("/update_ad", updatedb)
	//	meizu_router_upload.POST("/upload_behavior", upload_behavoir)
	//	meizu_router_upload.POST("/upload_behavior", upload_applist)
	//
	//	meizu_router_upload.POST("/update_ad_test", updatedb)
	//	meizu_router_upload.POST("/upload_behavior_test", upload_behavoir)
	//	meizu_router_upload.POST("/upload_applist_test", upload_applist)
	//}
	//
	//meizu_router_relevant := r.Group("/collaborate/meizu")
	//{
	//	meizu_router_relevant.POST("/search_relevant", relevant)
	//}

	r.Static("../templates","../templates")
	r.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/templates/test.html")
	})

	r.Run(":8000")
	//runtimeTicker.Stop()
}