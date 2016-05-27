package main

import (
	"bufio"
	"container/list"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	_ "odbc/driver"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

/**存放配置数据key-value形式*/
var CFTablesMap map[string]interface{} // {}为初始化成空
var baseTime time.Time

///////////////////////////////////////////////////

var reIDsMutex sync.Mutex
var reIDsClients list.List

var (
	redisCount = 4
)

func main() {

	tnow := time.Now()
	tmptime, _ := time.Parse("2006-01-02 15:04:05", "2015-06-19 17:01:47.590")
	subdur := tnow.Sub(tmptime)
	fmt.Printf("subdur: %f,%s ", subdur.Nanoseconds(), tmptime, " ")
	datastr := fmt.Sprintf("%d%d%d%d%d%d", tnow.Year(), tnow.Month(), tnow.Day(), tnow.Hour(), tnow.Minute(), tnow.Second())

	fmt.Println(time.Now().Format("2006-01-02_15:04:05"), "   ", datastr)
	//	test0()
	//	test1()rtrim(cast(data as CHAR(200))) as datacopy
	//	returndata := openDbString("select top 1 rtrim(cast(data as CHAR(200))) as datacopy,data from ATRes ")
	//	fmt.Println("result:", returndata)
	//openDbString("select code,data,CONVERT(CHAR(23), createtime, 121) as createtime,CONVERT(CHAR(23), updatetime, 121) as updatetime,groupCode,title,seq,valid,commentCount from ATRes ")
	CFTablesMap = loadConfig()
	//	test2()
	baseTime = time.Now()

	var N = 10
	sem := make(chan string, N)
	for i := 0; i < N; i++ {
		//		fmt.Println("index:", i)
		//		testLoginAndPost(i)
		go func(index int) {

			tnow := time.Now()
			starttimeint := (int)(tnow.Sub(baseTime).Seconds() * 1000000)

			fmt.Println("index:", index)

			///////////////////////////////////////
			TestMutexSpinLock()

			///////////////////////////////////////
			tmptime := time.Now()
			subdur := tmptime.Sub(tnow)
			tmpint := (int)(subdur.Seconds() * 1000)
			fmt.Printf("subdur:  ", tmpint)
			//			sem <- (fmt.Sprintf("%s,%d", tnow.Format("2006-01-02_15:04:05"), tmpint))
			//			sem <- (fmt.Sprintf("'%d:%d',%d", tnow.Minute(), tnow.Second(), tmpint))
			sem <- (fmt.Sprintf("%d,%d", starttimeint, tmpint))
		}(i)
	}
	//	var max = 0
	outString := "["
	for m := 0; m < N; m++ {
		//		<-sem
		tmp := <-sem
		//		if tmp > max {
		//			max = tmp
		//		}
		fmt.Println("FdMap:", tmp)
		//		outString += fmt.Sprintf("[%d,%d]", m%50, tmp)
		outString += fmt.Sprintf("[%s]", tmp)
		if m == N-1 {

		} else {
			outString += ","
		}
	}
	outString += "]"
	// fmt.Println("max:", outString)
	//	writeFileWithData(".config/data.json", outString, N)
	writeFileWithData("./config/output.html", outString, N)

}

// 封装测试数据
func getTableData(tablename string) string {
	var jsonstr string
	jsonstr += "\"table\":\"" + tablename + "\","
	jsonstr += "\"rows\":["

	jsonstr += "{" + getTableRowData(tablename) + "},"
	jsonstr += "{" + getTableRowData(tablename) + "}"

	jsonstr += "]"
	return jsonstr

}

// 封装测试数据
func getTableRowData(tablename string) string {
	var jsonstr string
	if strings.EqualFold(tablename, "ATResData") {
		jsonstr += "\"ResCode\":\"asd\","
		jsonstr += "\"comment\":\"3\","
		jsonstr += "\"seq\":\"1\","
		jsonstr += "\"valid\":\"1\","
		jsonstr += "\"usercode\":\"\","
		jsonstr += "\"createtime\":\"\","
		jsonstr += "\"guid\":\"" + GetGuid() + "\""
	} else if strings.EqualFold(tablename, "ATRes") {
		jsonstr += "\"data\":\"wbq.jpg\","
		jsonstr += "\"groupCode\":\"haha\","
		jsonstr += "\"title\":\"\","
		jsonstr += "\"valid\":\"1\","
		jsonstr += "\"seq\":\"1\","
		jsonstr += "\"commentCount\":\"4\","
		jsonstr += "\"code\":\"" + GetGuid() + "\""

	}
	return jsonstr
}

// 载入配置的json文件
func loadConfig() map[string]interface{} {

	CFTablesMap, err := readFile("./config/download.config")
	if err != nil {
		fmt.Println("readFile: ", err.Error())
		return nil
	}
	//	fmt.Println("map:", CFTablesMap["Tables"])
	//	tmpmap := CFTablesMap["Tables"].(map[string]interface{})
	//	fmt.Println("tmpmap:", tmpmap["ATResData"].(string))
	switch CFTablesMap["Tables"].(type) {
	case map[string]interface{}:
		//		tmpmap := CFTablesMap["Tables"].(map[string]interface{})
		//		fmt.Println("tmpmap:", tmpmap["ATResData"].(string))
		//		for k,v range tmpmap{}
	}
	return CFTablesMap
}

// 获取GUID唯一值
func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(b))) //使用zhifeiya名字做散列值，设定后不要变
	return hex.EncodeToString(h.Sum(nil))
	//    return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

// 将输入结果写入文件，data表示要写入的html文件内容，n用来命名文件头的 文件中存在%s用于写入数据代码
func writeFileWithData(filename string, data string, n int) {
	var tmpstring string
	f, _ := os.OpenFile(filename, os.O_RDONLY, 0666)
	defer f.Close()
	m := bufio.NewReader(f)
	char := 0
	words := 0
	lines := 0
	for {
		s, ok := m.ReadString('\n')
		//		fmt.Println(s)

		char += len(s)
		words += len(strings.Fields(s))
		lines++
		if ok != nil {
			break
		}
		if strings.Contains(s, "%s") {

			tmpstring += fmt.Sprintf(s, data) + "\n"
		} else {

			tmpstring += s + "\n"
		}
	}

	dirPath := "out"
	errdir := os.Mkdir(dirPath, 0)
	if errdir != nil {
		fmt.Println(errdir.Error())
	}
	//	tmptime, _ := time.Parse("2006-01-02_15:04:05", time.Now())
	fileName := dirPath + "/" + fmt.Sprintf("%d", n) + "-" + time.Now().Format("20060102_150405") + ".html"
	dstFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer dstFile.Close()
	dstFile.WriteString(tmpstring)
	fmt.Println("writeten to ", fileName)
}

// 读取json文件内容转换层map
func readFile(filename string) (map[string]interface{}, error) {
	FdMap := map[string]interface{}{}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return nil, err
	}
	if err := json.Unmarshal(bytes, &FdMap); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return nil, err
	}
	return FdMap, nil
}

// 字符串截取函数
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

///////////////////////////////////////////////////

func reserve_redis_client_pool() {

	for i := 0; i < redisCount; i++ {

		client, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("failed to create redis-client:", err.Error())
			os.Exit(1)
		}

		reIDsClients.PushBack(client)
	}

}

func TestMutexSpinLock() {

	client := pop_redis_client()
	if client == nil {
		fmt.Println("failed to pop new client")
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {

	}

	push_redis_client(client)
}

func pop_redis_client() *redis.Conn {
	reIDsMutex.Lock()
	defer reIDsMutex.Unlock()

	if reIDsClients.Len() == 0 {
		reserve_redis_client_pool()
	}

	conn := reIDsClients.Front().Value.(*redis.Conn)
	reIDsClients.Remove(reIDsClients.Front())

	return conn
}

func push_redis_client(client *redis.Conn) {
	reIDsMutex.Lock()
	defer reIDsMutex.Unlock()

	reIDsClients.PushBack(client)

	if redisCount > 100 && reIDsClients.Len() >= (redisCount*3/4) {
		count := redisCount / 2
		for i := 0; i < count && reIDsClients.Len() > 1; i++ {
			client := reIDsClients.Front().Value.(*redis.Conn)
			reIDsClients.Remove(reIDsClients.Front())
			(*client).Close()
		}

		redisCount -= count
	}
}
