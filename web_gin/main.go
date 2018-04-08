package main

import (
	"runtime"
	"github.com/gin-gonic/gin"
	"net/http"
	"io/ioutil"
	"bytes"
	"os/exec"
	"syscall"
	"os"
	"bufio"
	"strings"
	"github.com/spf13/cast"
)

type UrlQpsList []struct{
	Url string `json:"url"`
	Qps *string `json:"qps"`
}

type ReplayUrlQps struct{
	LatestTime string `json:"time"`
	HostPost string `json:"hostPort"`
	UrlQpsList []struct{
		Url string `json:"url"`
		Qps *string `json:"qps"`
	}`json:"urlList"`
}

const PATH = "/Users/cuixuange/goreplay"
const INPUTDIR = "/Users/cuixuange/goreplay/"
const StartBash = PATH + "/goreplay"
const confFile = "./conf"

func staticHtmlFunc(c *gin.Context){
	c.Redirect(http.StatusMovedPermanently, "./assets/templates/index.html")
}


// confMap => key:url value_list:qps + pgid + 是否被遍历过
func JudgeDiff(frontData UrlQpsList, logFileName string, outputIP string) (confMapFinal map[string][]string){
	confMap := make(map[string][]string,0)
	fin, err := os.OpenFile(confFile,os.O_RDONLY|os.O_CREATE,0644);if (err != nil){
		logs.Error("Fail to Open confFile,ERROR=%s",err)
		return
	}
	scanner := bufio.NewScanner(fin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		valueList := make([]string,3)
		line := strings.TrimSuffix(scanner.Text(), "\n")
		lineSplit :=  strings.Split(line,"\t")
		if lineSplit[2] != "NULL" {					//except logFile&IpPort just need url+qps
			valueList[0] = lineSplit[1] 			//qps
			valueList[1] = lineSplit[2] 			//pid
			valueList[2] = cast.ToString(false)
			confMap[lineSplit[0]] = valueList
		}
	}
	for _,front := range(frontData){
		if value,ok := confMap[front.Url];ok{
			if *front.Qps != value[0]{
				//kill
				StopGoreplay(cast.ToInt(value[1]))
				//start
				pgidNew := StartGoreplay(front.Url,value[0],logFileName,outputIP)
				//update confMap
				confMap[front.Url][0] = *front.Qps
				confMap[front.Url][1] = cast.ToString(pgidNew)
				confMap[front.Url][2] = cast.ToString(true)
			}else{
				//url qps equal so falg ==true
				confMap[front.Url][2] = cast.ToString(true)
			}
		}else{
			//start; when url not in confMap
			pgidNew := StartGoreplay(front.Url,*front.Qps,logFileName,outputIP)
			//update confMap
			valueList := make([]string,3)
			valueList[0] = *front.Qps
			valueList[1] = cast.ToString(pgidNew)
			valueList[2] = cast.ToString(true)
			confMap[front.Url] = valueList
		}
		logs.Info("after start: ",confMap)
	}

	// 删除confMap中没有被遍历到的数据   直接删除iterator会失效么？
	// kill 对应pgid
	for key := range confMap{
		if len(confMap) != 0 && confMap[key][2] == "false" && confMap[key][1]!="NULL"{
			//kill
			logs.Info("kill false",confMap)
			StopGoreplay(cast.ToInt(confMap[key][1]))
			delete(confMap,key)
		}
	}
	return confMap
}

// key:url value:qps pgid bool
//last line is last logfile
func WriteToFile(finalRes map[string][]string,logFileName string,outputIP string){
	fo, err := os.OpenFile(confFile,os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND,0644);if (err != nil){
		logs.Error("Fail to Write confFile,ERROR=%s",err)
		return
	}
	defer fo.Close()
	for key := range finalRes{
		var tmp = key+"\t"+finalRes[key][0]+"\t"+finalRes[key][1]
		fo.Write([]byte(tmp + "\n"))
	}
	fo.Write([]byte("lastLogFile"+"\t" + logFileName + "\t"+"NULL"+"\n"))
	fo.Write([]byte("lastHostPost"+"\t" + outputIP + "\t"+"NULL"+"\n"))
}

func StopGoreplay(pgid int){
	//need root 权限  或者sudo 启动go
	err := syscall.Kill(-pgid, syscall.SIGKILL)				//kill all childs
	if err != nil{
		logs.Error("process kill failed , PGID=%d ,ERROR=%s",pgid,err)
	}else{
		logs.Info("process kill Success , PGID=%d",pgid)
	}

}
func StartGoreplay(URL string, qps string, logFileName string, outputIP string) (PGID int){
	inputFile := " --input-file " + "'" + INPUTDIR + logFileName + "|" + qps + "'"
	allowURL := " --http-allow-url " + URL
	outputPath := " --output-http "+ outputIP+" --input-file-loop"
	//logs.Info("inputfile: ",inputfile)
	handleBash := StartBash + inputFile + allowURL + outputPath
	logs.Info("bash: ",handleBash)
	cmd := exec.Command("/bin/sh", "-c", "sudo " + handleBash)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := cmd.Start()
	var out bytes.Buffer
	cmd.Stdout = &out
	if err != nil {
		logs.Error("",err)
	}
	//logs.Info("/bin/sh: PID=%d , RES=%q", cmd.Process.Pid, out.String())			//应该是PPID？？
	pgid , _ := syscall.Getpgid(cmd.Process.Pid)
	logs.Info("/bin/sh: pgid=%d , RES=%q", pgid, out.String())
	return pgid
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
		logFileName  := req.LatestTime				//init input log file
		outputIP := req.HostPost					//init output hostPost
		for _,value := range(req.UrlQpsList){
			logs.Info("url:%s qps:%s",value.Url,*value.Qps)
			strQps := "100%"
			if value.Qps == nil{
				value.Qps = &strQps
			}
		}
		confMapFinal := JudgeDiff(req.UrlQpsList,logFileName,outputIP)
		logs.Info("confMapFinal: ",confMapFinal)
		WriteToFile(confMapFinal,logFileName,outputIP)
	}
}

func reloadConfByJsonList(c *gin.Context){
	confMap := make(map[string]string,0)
	fin, err := os.OpenFile(confFile,os.O_RDONLY,0644);if (err != nil){
		logs.Error("Fail to Open confFile,ERROR=%s",err)
		return
	}
	scanner := bufio.NewScanner(fin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\n")
		lineSplit :=  strings.Split(line,"\t")
		//logs.Error("",len(lineSplit),lineSplit)
		confMap[lineSplit[0]] = lineSplit[1]			//url:qps and lastLogfile:name
	}
	//logs.Info("",confMap)
	c.JSON(200,confMap)
}

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	gor := gin.Default()
	gor.Static("/assets", "./assets")
	gor.GET("/index",staticHtmlFunc)

	gor.POST("/submitAndStore",GetJsonListFromPost)
	gor.POST("/reload",reloadConfByJsonList)
	gor.Run(":6790")
}