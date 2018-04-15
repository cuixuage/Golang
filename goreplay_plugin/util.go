package main

import (
	"os"
	"bufio"
	"strings"
	"syscall"
	"github.com/spf13/cast"
	"io/ioutil"
	"time"
	"sort"
)
//alive goReplay process
type UrlQpsPid struct{
	Url string
	Qps string
	Pid string
	LogFile string
	HostIP string
}

//restart notAlive process
func (value *UrlQpsPid)KeepProAlive() {
	killErr := syscall.Kill(cast.ToInt(value.Pid), syscall.Signal(0))
	if killErr != nil {
		logs.Info("dead and restart",killErr)
		value.Restart()
	}
}


func (value *UrlQpsPid) StopGoReplay() {
	err := syscall.Kill(-cast.ToInt(value.Pid), syscall.SIGKILL)				//kill all childs
	if err != nil{
		logs.Error("process kill failed , PGID=%s ,ERROR=%s",value.Pid,err)
	}else{
		logs.Info("process kill Success , PGID=%s",value.Pid)
	}
}

//before start need be dead
func (value *UrlQpsPid) Restart(){
	newPid := StartNewGoReplay(value.Url,value.Qps,value.LogFile,value.HostIP)
	value.Pid = cast.ToString(newPid)
}

func (value *UrlQpsPid) RestartByQps(newQps string){
	newPid := StartNewGoReplay(value.Url,newQps,value.LogFile,value.HostIP)
	value.Qps = newQps
	value.Pid = cast.ToString(newPid)
}

func updateSecondLogFile(){
	logs.Info("call updateSecondLogFile=%s", LatestLogFile)
	files, err1 := ioutil.ReadDir(InputFileDir)
	if err1 != nil {
		logs.Error("", err1)
		return
	}
	type sortItem  struct {
		timeSub time.Duration
		fileName string
	}
	sortSlice := make([]sortItem,len(files))
	var fileName string
	timeNow := time.Now()
	for key,f := range files {
		sortSlice[key].timeSub = timeNow.Sub(f.ModTime())
		sortSlice[key].fileName = f.Name()
	}
	sort.Slice(sortSlice[:], func(i, j int) bool {
		return sortSlice[i].timeSub < sortSlice[j].timeSub
	})
	if len(sortSlice)<=1{
		fileName = sortSlice[0].fileName
	}else{
		fileName = sortSlice[1].fileName
	}
	//if fileName != LatestLogFile && StrictUpdatePro(LatestLogFile,sortSlice){
	if fileName != LatestLogFile{
		LatestLogFile = fileName
		logs.Info("updateFile=%s", LatestLogFile)
		confMap := ReadConfLogFile()
		confMap["lastLogFile"] = LatestLogFile
		WriteAllToLogFile(confMap)			//firstly update logFileName
		RestartAllPro()
	}
}

//read url qps with pid
func InitAliveList() (listAlive []UrlQpsPid){
	confMapLog := ReadConfLogFile()
	logFileName := confMapLog["lastLogFile"]
	outputIP := confMapLog["lastHostPost"]
	listAlives := make([]UrlQpsPid,0)
	fin, err := os.OpenFile(confGoReplay,os.O_RDONLY|os.O_CREATE,0644);if (err != nil){
		logs.Error("Fail to Open confFile,ERROR=",err)
		return
	}
	defer fin.Close()
	scanner := bufio.NewScanner(fin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		oneUrl := new (UrlQpsPid)
		line := strings.TrimSuffix(scanner.Text(), "\n")
		lineSplit :=  strings.Split(line,"\t")
		oneUrl.Url = lineSplit[0]
		oneUrl.Qps = lineSplit[1]
		oneUrl.Pid = lineSplit[2]
		oneUrl.LogFile = logFileName
		oneUrl.HostIP = outputIP
		listAlives = append(listAlives,*oneUrl)
	}
	return listAlives
}

//read url qps without pid in map
func ReadConfGoRepaly()(rsp map[string]string){
	confMap := make(map[string]string,0)
	fin, err := os.OpenFile(confGoReplay,os.O_RDONLY|os.O_CREATE,0644);if (err != nil){
		logs.Error("Failed ReadConfGoRepaly,ERROR=",err)
		return
	}
	defer fin.Close()
	scanner := bufio.NewScanner(fin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\n")
		lineSplit :=  strings.Split(line,"\t")
		confMap[lineSplit[0]] = lineSplit[1]
	}
	return confMap
}

//read current logFile and HostIP
func ReadConfLogFile()(rsp map[string]string){
	confMap := make(map[string]string,0)
	fin, err := os.OpenFile(confFile,os.O_RDONLY|os.O_CREATE,0644);if (err != nil){
		logs.Error("Failed ReadConfLogFile,ERROR=",err)
		return
	}
	defer fin.Close()
	scanner := bufio.NewScanner(fin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\n")
		lineSplit :=  strings.Split(line,"\t")
		confMap[lineSplit[0]] = lineSplit[1]
	}
	return confMap
}

func WriteToReplay(listAlivePro []UrlQpsPid){
	fo, err := os.OpenFile(confGoReplay,os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND,0644);if (err != nil){
		logs.Error("Fail to Write confFile,ERROR=",err)
		return
	}
	defer fo.Close()
	for _,value := range listAlivePro{
		var tmp = value.Url+"\t"+value.Qps+"\t"+value.Pid
		fo.Write([]byte(tmp + "\n"))
	}
}

func WriteToLogFile(logFileName string,outputIP string){
	confMap := ReadConfLogFile()
	confMap["lastLogFile"] = logFileName
	confMap["lastHostPost"] = outputIP
	WriteAllToLogFile(confMap)
}

func WriteAllToLogFile(confMap map[string]string){
	fo, err := os.OpenFile(confFile,os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND,0644);if (err != nil){
		logs.Error("Fail to Write confFile,ERROR=",err)
		return
	}
	defer fo.Close()
	for key,value := range confMap{
		token := key+"\t"+value+"\n"
		fo.Write([]byte(token))
	}
}


func GetPortFromConf(filePath string)(string){
	fin, err := os.OpenFile(filePath,os.O_RDONLY|os.O_CREATE,0644);if (err != nil){
		logs.Error("Failed ReadConfLogFile,ERROR=",err)
		return ""
	}
	defer fin.Close()
	scanner := bufio.NewScanner(fin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\n")
		lineSplit :=  strings.Split(line,"=")
		return lineSplit[1]
	}
	return ""
}