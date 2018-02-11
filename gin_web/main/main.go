package main

import (
	"runtime"
	//"code.byted.org/gin/ginex"
	"net/http"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"code.byted.org/gin/ginex/ctx"
	"github.com/spf13/cast"
	"reflect"
	"meizuapi/toutiao_meizu_app"
	"protobuf/proto"
)

type recoverWriter struct{}
func (rw *recoverWriter) Write(p []byte) (int, error) {
	return gin.DefaultErrorWriter.Write(p)
}
func JsonRequestContext() gin.HandlerFunc {
		return func(c *gin.Context) {
			if  c.Request.Method != "GET" {
				cacheBody, err := ioutil.ReadAll(c.Request.Body)              //打开
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(cacheBody)) //无操作关闭
					if err == nil && cacheBody != nil{
						c.Set("cachebody", cacheBody)
					}

				checkreq := &CheckRequst{}
				err = json.Unmarshal([]byte(c.Request.Header.Get("METADATA")), &checkreq)
				_,e := get_check_error_response(checkreq)
				if e != nil{
					var res *JsonResponse = &JsonResponse{}
					res.Code = cast.ToInt(META_DATA_ERROR)
					res.Message = "catch error 校验失败"
					c.AbortWithStatusJSON(500,res)
					c.Next()
				}
				c.Next()
		}
	}
}

func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		var res *JsonResponse = &JsonResponse{}
		if c.IsAborted(){
			res.Code = 500
			res.Message = "catch error"
			c.AbortWithStatusJSON(500,res)
		}else if c.Request.Method != "GET" {
			pb,exits := c.Get("pb")
			if exits && pb=="false" {
				res.Code = cast.ToInt(*SUCCESS)
				res.Message = "success"
				res.Value = map[string]int{"status":0}
				res.Redirect = ""
				c.JSON(200,res)
			}else if pb == "upload"{
				strReqId,_ := c.Get("uploadReqId")
				rspUpload := &JsonResponseUpload{}
				rspUpload.ReqId = strReqId.(string)
				rspUpload.Message = StrSuccess
				rspUpload.Code = cast.ToInt(SUCCESS)
				c.JSON(200,rspUpload)
			} else if pb == "true"{
				rsp,_ := c.Get("rsp")
				pb_result := &toutiao_meizu_app.BinaryDisplayResult{}
				pb_result.Code = SUCCESS
				pb_result.Message = &StrSuccess
				r := reflect.ValueOf(rsp)
				f_rankresult := reflect.Indirect(r).FieldByName("RankResults")
				pb_result.Value = f_rankresult.Bytes()
				pb_bytes,_ := proto.Marshal(pb_result)
				c.Data(200,"application/octet-stream",pb_bytes)
			}
		}
	}
}


func main(){
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	r := gin.Default()

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
		meizu_router.POST("/idea_predict/", idea_predict)
		meizu_router.POST("/recommend/normal", recommend_predict)
		meizu_router.POST("/predict/callback/", callback)
		meizu_router.POST("/predict_binary/", predict_binary)
		/*
		FIXME  路径是search_pb   ??  search路径是按照json格式返回的
 		*/
		meizu_router.POST("/search_pb/", search_pb)
		meizu_router.POST("/upload_behavior/", upload_behavior)
		meizu_router.POST("/upload_applist/", upload_applist)

		meizu_router.POST("/idea_predict_test/", idea_predict_test)
		meizu_router.POST("/recommend/normal_test", recommend_predict_test)
		meizu_router.POST("/predict/callback_test/", callback_test)
		meizu_router.POST("/predict_binary_test/", predict_binary_test)
		meizu_router.POST("/search_pb_test/", search_pb_test)
		meizu_router.POST("/upload_behavior_test/", upload_behavior_test)
		meizu_router.POST("/upload_applist_test/", upload_applist_test)


		//这个updata——db路径在使用吗/
		//url(r'^collaborate/meizu/update_ad/$', 'updatedb'),
	}

	r.Static("../templates","../templates")
	r.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/templates/test.html")
	})

	r.Run(":8000")
	//runtimeTicker.Stop()
}