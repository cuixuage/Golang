package main
import (
	"net"
	"meizuapi/thrift_gen/predict"
	"code.byted.org/gopkg/thrift"
	"strings"
	"math/rand"
	"code.byted.org/kite/kitc"
	//"code.byted.org/gopkg/logs"
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
const TIMEOUT_CONNECT_TIME = 300
//const RETRY_CONNECT_COUNT = 3

var pool connpool.ConnPool
//var predictConf map[string]string
//var connTagMutex *sync.RWMutex

func newInstances(hosts []string) []*kitc.Instance {
	var ins []*kitc.Instance
	for _, hostPort := range hosts {
		val := strings.Split(hostPort, ":")
		if len(val) == 2 {
			ins = append(ins, kitc.NewInstance(val[0], val[1], make(map[string]string)))
		}
	}
	return ins
}

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
func getConn(xs []*kitc.Instance) (c net.Conn, err error) {
	for _, i := range xs {
		c, err = pool.Get(i.Host(), i.Port(), TIMEOUT_CONNECT_TIME*time.Millisecond) //FIXME 注意下时间
		//fmt.Print("RemoteConn:")
		//fmt.Print(c.RemoteAddr().String()," ","local:")
		//fmt.Print(c.LocalAddr().String())
		//fmt.Print("\n")
		return c,err
	}
	return
}

func getPredictClient(conn net.Conn) *predict.PredictClient {
	var transport thrift.TTransport
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTBufferedTransportFactory(4096))
	transport = thrift.NewTSocketFromConnTimeout(conn, TIMEOUT_CONNECT_TIME * time.Millisecond)
	transport = transportFactory.GetTransport(transport)
	return predict.NewPredictClientFactory(transport, protocolFactory)
}
func checkResponse(rsp interface{}, err error) (interface{}, error) {
	if err != nil || rsp==nil  {
		return rsp,err
	}
	r := reflect.ValueOf(rsp)
	f_status := reflect.Indirect(r).FieldByName("Status")
	status := f_status.String()

	if status != "success" {
		//msg := reflect.Indirect(r).FieldByName("ErrMessage").Elem().String()
		//return rsp, errors.New(fmt.Sprintf("Predict error %d %s", status, msg))
		return rsp, errors.New(fmt.Sprintf("Predict error %d %s", status, "status error"))
	}
	fmt.Print("\n")
	fmt.Print("checkresponse status is:",status,"\n")
	fmt.Print("\n")
	return rsp, nil
}
func CallClient(name string, req interface{}) (interface{}, error) {

	var ins []*kitc.Instance
	ins = randomInstances(AllServers)
	conn, err := getConn(ins)
	//defer pool.Put(conn, err)
	//if err != nil {
	//	logs.Error("connect all err", err)
	//	return nil, err
	//}
	rsp, err := doCallClient(name, req, conn)			//响应
	if err != nil {
		fmt.Printf("call fail %s \n %s \n conn=%s \n req=%s \n rsp=%s\n", name, err, conn.RemoteAddr(), req, rsp)
	}
	return rsp, err
}

func doCallClient(name string, req interface{}, conn net.Conn) (interface{}, error) {
	var client *predict.PredictClient
	client = getPredictClient(conn)

	switch name {

	case "callback":
		return checkResponse(client.UploadServerImprOneway(req.(*predict.Req)),nil)

	case "idea_predict":
		return checkResponse(client.FeedPredict(req.(*predict.Req)))

	case "predictbinary":
		//fmt.Print("\n")
		//fmt.Print("this is doCallClient req\n")
		//fmt.Print(req.(*predict.Req))
		//fmt.Print("\n")
		return checkResponse(client.BinaryPredict(req.(*predict.Req)))

	case "recommend":
		//fmt.Print("\n")
		//fmt.Print("this is doCallClient req\n")
		//fmt.Print(req.(*predict.Req))
		//fmt.Print("\n")
		return checkResponse(client.RelatePredict(req.(*predict.Req)))

	case "search":
		return checkResponse(client.BinarySearch(req.(*predict.Req)))

	case "upload_behavior":
		//fmt.Print("\n")
		//fmt.Print("this is doCallClient req\n")
		//fmt.Print(req.(*predict.AckReq))
		//fmt.Print("\n")
		return checkResponse(client.Ack(req.(*predict.AckReq)))

	case "upload_applist":
		fmt.Print("\n")
		fmt.Print("this is doCallClient req\n")
		fmt.Print(req.(*predict.AckReq))
		fmt.Print("\n")
		return checkResponse(client.UploadApplist(req.(*predict.AckReq)))

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
	return hostPort
}
func InitThriftClient() {

	 //hosts := splitLine()
	hosts := []string{"10.8.64.232:7270"}
	AllServers = newInstances(hosts)
	pool = connpool.NewShortPool()
	//fmt.Print("init thrift ok\n")
}
