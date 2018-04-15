package main

import (
	"runtime"
	"github.com/gin-gonic/gin"
	"net/http"
	"io/ioutil"
	"bytes"
	"os/exec"
	"syscall"
	"github.com/spf13/cast"
	"time"
)

var LatestLogFile string
var StartPath string
var InputFileDir string
var TickerTime int64
var ticker *time.Ticker
const confGoReplay = "./conf/confReplay"
const confFile = "./conf/confLogfile"


type UrlQpsList []struct{
	Url string `json:"url"`
	Qps *string `json:"qps"`
}

type ReplayUrlQps struct{
	LatestTimeFile string `json:"newLogFile"`
	HostPost string `json:"hostPort"`
	StartPath string `json:"startPath"`
	InputFileDir string `json:"inputDir"`
	TickerTime string `json:"tickerTime"`
	UrlQpsList []struct{
		Url string `json:"url"`
		Qps *string `json:"qps"`
	}`json:"urlList"`
}

func staticHtmlFunc(c *gin.Context){
	c.Redirect(http.StatusMovedPermanently, "./assets/templates/index.html")
}

func StartNewGoReplay(URL string, qps string, logFileName string, outputIP string) (PGID int){
	inputFile := " --input-file " + "'" + InputFileDir + logFileName + "'"
	allowURL := " --http-allow-url " + URL
	outputPath := " --output-http "+ "'" + outputIP + "|" + qps + "'" + " --input-file-loop"
	handleBash := StartPath + inputFile + allowURL + outputPath
	logs.Info("bash: ",handleBash)
	cmd := exec.Command("/bin/sh", "-c", "sudo " + handleBash)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := cmd.Start()
	var out bytes.Buffer
	cmd.Stdout = &out
	if err != nil {
		logs.Error("",err)
	}
	pgid , _ := syscall.Getpgid(cmd.Process.Pid)
	logs.Info("/bin/sh: pgid=%d , RES=%q", pgid, out.String())
	return pgid
}

func StoreBeforeOrNot (url string,qps string)(int){
	aliveList := InitAliveList()
	for _,value := range aliveList {
		if value.Url==url && value.Qps==qps{
			return 0
		}
		if value.Url==url && value.Qps!=qps{
			return 1		//need update
		}
	}
	return -1			//need start newProcess
}

func RestartFromAlive(url string,newQps string,aliveList []UrlQpsPid)([]UrlQpsPid){
	for key,value := range aliveList {
		if value.Url==url {
			//kill origin
			value.StopGoReplay()
			//start new goreplay
			value.RestartByQps(newQps)
			//update list
			aliveList[key] = value
			break
		}
	}
	return aliveList
}

func StartToAlive(url string,qps string,aliveList[]UrlQpsPid)(res []UrlQpsPid){
	confMapLog := ReadConfLogFile()
	logFileName := confMapLog["lastLogFile"]
	outputIP := confMapLog["lastHostPost"]
	element := new (UrlQpsPid)
	newPid := StartNewGoReplay(url,qps,logFileName,outputIP)
	element.Url = url
	element.Qps = qps
	element.Pid = cast.ToString(newPid)
	aliveList = append(aliveList,*element)
	return aliveList
}

func DelItemFromAliveList(delUrl string,aliveList []UrlQpsPid) (res []UrlQpsPid){
	for i,v := range aliveList{
		if v.Url == delUrl{
			v.StopGoReplay()
			aliveList[i] = aliveList[len(aliveList)-1] // Replace it with the last one
			aliveList = aliveList[:len(aliveList)-1]
			break
		}
	}
	return aliveList
}

func JudgeDiff(frontData UrlQpsList) (finalAlive []UrlQpsPid){
	//keep listAlivePro the same
	listAlivePro := InitAliveList()
	mapAlivePro := make(map[string]bool,0)
	for _,value := range listAlivePro {
		mapAlivePro[value.Url] = false				//init 都没有遍历
	}

	for _,front := range frontData {
		var res int = StoreBeforeOrNot(front.Url, *front.Qps)
		if res == 1 {
			listAlivePro = RestartFromAlive(front.Url, *front.Qps,listAlivePro)
		}else if res == -1{
			listAlivePro = StartToAlive(front.Url, *front.Qps,listAlivePro)
		}
		mapAlivePro[front.Url] = true
	}
	//kill conf origin没有遍历到的process
	for key,value := range mapAlivePro{
		if value == false{
			listAlivePro = DelItemFromAliveList(key,listAlivePro)
		}
	}
	return listAlivePro
}

func UpdateReplayConf(req ReplayUrlQps)(map[string]string){
	confMap := ReadConfLogFile()
	confMap["lastLogFile"] = req.LatestTimeFile
	confMap["lastHostPost"] = req.HostPost
	confMap["inputFileDir"] = req.InputFileDir
	confMap["startPath"] = req.StartPath
	confMap["tickerTime"] = req.TickerTime
	StartPath = confMap["startPath"]
	return confMap
}
func GetJsonListFromPost(c *gin.Context){
	cacheBody, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(cacheBody))
	//logs.Info("%s header=%s body= %s", c.Request.URL.Path, jsonMeta, string(cacheBody))
	req := &ReplayUrlQps{}
	err := c.BindJSON(req)
	if err != nil{
		logs.Error("bind json error=",err)
		c.Abort()
		return
	}else{
		confMap := UpdateReplayConf(*req)
		WriteAllToLogFile(confMap)

		if InputFileDir != req.InputFileDir && LatestLogFile != req.LatestTimeFile{		//first start
			InputFileDir = req.InputFileDir
			LatestLogFile = req.LatestTimeFile
			RestartAllPro()
		} else if  InputFileDir != req.InputFileDir{				//update inputDir
			InputFileDir = req.InputFileDir
			RestartAllPro()
		} else if  LatestLogFile != req.LatestTimeFile{			//updateLogFile
			LatestLogFile = req.LatestTimeFile
			RestartAllPro()
		}
		if TickerTime != cast.ToInt64(req.TickerTime){				//update timeTicker
			TickerTime = cast.ToInt64(req.TickerTime)
			if ticker != nil{
				ticker.Stop()
			}
			ticker = time.NewTicker(time.Second * time.Duration(TickerTime))
			logs.Info("new ticker",time.Second * time.Duration(TickerTime))
			go ScanLatestFile(ticker)
		}
		for key,value := range req.UrlQpsList{
			//logs.Info("url:%s qps:%s",value.Url,*value.Qps)
			strQps := "0"
			if value.Qps == nil{
				req.UrlQpsList[key].Qps = &strQps
			}
		}

		listAlivePro := JudgeDiff(req.UrlQpsList)
		logs.Info("UpdateFromSubmit: ",listAlivePro)
		WriteToReplay(listAlivePro)
	}
}

func RestartAllPro(){
	//store all url+qps
	aliveProList := InitAliveList()
	for key,value := range aliveProList {
		//kill origin
		value.StopGoReplay()
		//start new goreplay
		value.Restart()
		//update list
		aliveProList[key].Pid = cast.ToString(value.Pid)
	}
	//update conf pid
	logs.Info("listAlivePro: ",aliveProList)
	WriteToReplay(aliveProList)
}

func ReloadConfByJsonList(c *gin.Context){
	confMap := ReadConfGoRepaly()
	confMapLog := ReadConfLogFile()
	for key,value := range confMap{
		confMapLog[key] = value
	}
	c.JSON(200,confMapLog)
}

//func StrictUpdatePro(curLogFile string, allLogFile []sortItem) (bool){
//	currentKey := 0
//	for key,value := range allLogFile{
//		if value.fileName == curLogFile{
//			currentKey = key
//			break
//		}
//	}
//	if currentKey >= len(allLogFile)/2{
//		logs.Info("currentKey=", currentKey)
//		return true
//	}
//	return false
//}


func ScanLatestFile(ticker *time.Ticker) {
	for range ticker.C {
		flag := false
		aliveList := InitAliveList()
		for i,value := range aliveList{
			value.KeepProAlive()	 //updatePid
			//logs.Info("cur pid=%d,new pid=%d",aliveList[i].Pid,value.Pid)
			if aliveList[i].Pid != value.Pid{
				aliveList[i].Pid = value.Pid
				flag = true
			}
		}
		if len(aliveList)!=0 && flag == true{
			logs.Info("Scan to update aliveList:",aliveList)
			WriteToReplay(aliveList)
		}
		//avoid start deadProcess twice when go restart
		updateSecondLogFile()
	}
}

func initGlobal(){
	confMap := ReadConfLogFile()
	if val,ok := confMap["lastLogFile"];ok{
		LatestLogFile = val
	}
	if val,ok := confMap["inputFileDir"];ok{
		InputFileDir = val
	}
	if val,ok := confMap["startPath"];ok{
		StartPath = val
	}
	if val, ok := confMap["tickerTime"]; ok {
		TickerTime = cast.ToInt64(val)
		ticker = time.NewTicker(time.Second * time.Duration(cast.ToInt64(val)))
		go ScanLatestFile(ticker)
	}
}
func main(){
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	gin.SetMode(gin.ReleaseMode)
	initGlobal()
	gor := gin.Default()
	gor.Static("/assets", "./assets")
	gor.GET("/",staticHtmlFunc)

	gor.POST("/submitAndStore",GetJsonListFromPost)
	gor.POST("/reload",ReloadConfByJsonList)

	port := GetPortFromConf("./conf/confPort")
	if port == ""{
		logs.Error("need start port")
		runtime.Goexit()
	}
	gor.Run("0.0.0.0:"+port)
}