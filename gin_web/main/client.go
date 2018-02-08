package main
import (
	"net"
	"meizuapi/thrift_gen/predict"
	"code.byted.org/gopkg/thrift"
	"strings"
	"math/rand"
	"code.byted.org/kite/kitc"
	"code.byted.org/gopkg/logs"
	"sync"
	"time"
	"errors"
	"fmt"
	"code.byted.org/kite/kitc/connpool"
	"reflect"
	"os"
	"bufio"
	"unicode"
	"io"
)

var AllServers []*kitc.Instance				//在initthrift 中初始化
//const HASH_WINDOW_SIZE = 10
//const TAG_LAST_ERROR_KEY = "_error_time"
//const BLOCK_TIMEOUT_CONNECT_TIME = 5
//const RETRY_CONNECT_COUNT = 3

var pool connpool.ConnPool
var predictConf map[string]string
var connTagMutex *sync.RWMutex

func newInstances(hosts []string) []*kitc.Instance {
	var ins []*kitc.Instance
	for _, hostPort := range hosts {
		val := strings.Split(hostPort, ":")
		if len(val) == 2 {
			//fmt.Print("\n")
			//fmt.Print(val[0],"   ", val[1])
			//fmt.Print("\n")
			ins = append(ins, kitc.NewInstance(val[0], val[1], make(map[string]string)))
		}
	}
	return ins
}

//重新排序instance
func randomInstances(xs []*kitc.Instance) []*kitc.Instance {
	rand.Seed(time.Now().UnixNano())
	randXS := make([]*kitc.Instance, len(xs))
	copyXS := make([]*kitc.Instance, len(xs))
	copy(copyXS, xs)
	copy(randXS, xs)
	n := len(xs)
	for i := range copyXS {
		rand_idx := rand.Intn(n)
		randXS[i] = copyXS[rand_idx]
		copyXS[rand_idx] = copyXS[n-1]
		n--
	}
	return randXS
}
/*
1。从随机顺序的server 加读锁 判断tags情况(instance是否可以链接?)
2. 连续三次没有成功连接到server直接退出
 */
func getConn(xs []*kitc.Instance) (c net.Conn, err error) {
	for _, i := range xs {
		c, err = pool.Get(i.Host(), i.Port(), 300*time.Millisecond) //FIXME 注意下时间
		fmt.Print("this is getconn func \n")
		fmt.Print(c.RemoteAddr().String()," ","local:")
		fmt.Print(c.LocalAddr().String())
		fmt.Print("\n")
		return c,err
	}
	return
}

func getPredictClient(conn net.Conn) *predict.PredictClient {
	var transport thrift.TTransport
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTBufferedTransportFactory(4096))
	transport = thrift.NewTSocketFromConnTimeout(conn, 300 * time.Millisecond)
	//fmt.Print("\n")
	//fmt.Print(conn.RemoteAddr().String())
	//fmt.Print("\n")
	transport = transportFactory.GetTransport(transport)
	return predict.NewPredictClientFactory(transport, protocolFactory)
}
/*
1. thrift返回的err!=nil 使得对应handler abort
2. err==nil 同时状态码==success
 */
func checkResponse(rsp interface{}, err error) (interface{}, error) {
	if err != nil || rsp==nil  {
		return rsp,err
	}
	r := reflect.ValueOf(rsp)
	f_status := reflect.Indirect(r).FieldByName("Status")
	f_rankresult := reflect.Indirect(r).FieldByName("RankResults")
	status := f_status.String()
	//version := f_version.String()
	if status != "success" {
		msg := reflect.Indirect(r).FieldByName("ErrMessage").Elem().String()
		return rsp, errors.New(fmt.Sprintf("Predict error %d %s", status, msg))
	}
	fmt.Print("\n")
	fmt.Print("checkresponse status is:",status,"\n")
	fmt.Print(string(f_rankresult.Bytes()))
	fmt.Print("\n")

	return rsp, nil
}
/*
1。 random序列的server slice
2。 getconn 寻找可链接的server
3。 已构造request doCallClient(req coon)   链接对应thrift得到response
4。 返回 server thrift rsp
*/
func CallClient(name string, req interface{}) (interface{}, error) {

	var ins []*kitc.Instance
	ins = randomInstances(AllServers)
	conn, err := getConn(ins)
	//defer pool.Put(conn, err)
	if err != nil {
		logs.Error("connect all err", err)
		return nil, err
	}
	rsp, err := doCallClient(name, req, conn)			//响应
	if err != nil {
		fmt.Printf("call fail %s %s conn=%s req=%s rsp=%s", name, err, conn.RemoteAddr(), req, rsp)
	}
	return rsp, err
}

func doCallClient(name string, req interface{}, conn net.Conn) (interface{}, error) {
	var client *predict.PredictClient
	client = getPredictClient(conn)

	//fmt.Print("\n")
	//test := req.(*predict.Req)
	//fmt.Print("docallclient"," = ",test.Debug," ")
	//fmt.Print("docallclient"," = ",test.Imei)
	//fmt.Print("\n")

	switch name {
	//case "binary_predict":
	//return client.BinaryPredict(req.(*meizu.PredictBinaryPredictArgs))
	case "callback":
		return checkResponse(client.UploadServerImprOneway(req.(*predict.Req)),nil)

	case "idea_predict":
		return checkResponse(client.FeedPredict(req.(*predict.Req)))

	case "predictbinary":
		return checkResponse(client.BinaryPredict(req.(*predict.Req)))
	//
	//case "Predict":
	//	return checkResponse(client.Predict(req.(*vivo.PredictReq)))
	//
	//case "Search":
	//	return checkResponse(client.Search(req.(*vivo.PredictReq)))
	//
	//case "Upload":
	//	return checkResponse(client.Upload(req.(*vivo.BehaviorReq)))
	//
	//case "Related":
	//	return checkResponse(client.Related(req.(*vivo.PredictReq)))

	}
	err := errors.New("invalid method name")
	return nil, err
}
func c() string{
	return ""
}
func getWholeLine(scanner *bufio.Scanner) (line string, ok bool) {
	ok = false
	var whole_line string
	for scanner.Scan() {
		ok = true
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, "\\") {
			// merge line ends with backslash
			line = line[:len(line)-1]
			whole_line += line
			continue
		}
		whole_line += line
		break
	}
	return whole_line, ok
}
func isSsConfSep(r rune) bool {
	if r == ':' || r == '=' || unicode.IsSpace(r) {
		return true
	}
	return false
}

func splitLine()  []string {
	var serverHost []string
	var serverPort []string
	conf := "/Users/cuixuange/go/src/2018_2_5/conf/thriftServer.conf"
	file, _ := os.Open(conf)
	defer file.Close()
	scanner := bufio.NewReader(file)
	for {
		line, _, c := scanner.ReadLine()
		if c == io.EOF {
			break
		}
		result := strings.Split(string(line), " ")
		//fmt.Print(result[0])
		if result[0] == "meizu_ad_predict_host" {
			serverHost = result
		}
		if result[0] == "meizu_ad_predict_port" {
			serverPort = result
		}
	}
	var hostPort = make([]string,len(serverHost))
	for i,_ := range serverHost{
		if i >= 1 {
			hostPort[i] = serverHost[i] + ":" + serverPort[i]
		}
	}
	for i,_ := range hostPort{
		hostPort[i] = strings.Replace(hostPort[i], ",", "", -1)
	}
	//for _,val := range hostPort{
	//	fmt.Print(val," ")
	//}
	return hostPort
}
func InitThriftClient() {
	/*
	 mutex Allservers pool初始化
	 */
	hosts := splitLine()
	//hosts :=[]string {"10.3.29.13:7250","10.3.29.14:7250"}
	connTagMutex = &sync.RWMutex{}
	AllServers = newInstances(hosts)
	pool = connpool.NewShortPool()
	fmt.Print("init thrift ok\n")
	//InitABConfig()
}

//func main(){
//	parseConf()
//	InitThriftClient()
//}