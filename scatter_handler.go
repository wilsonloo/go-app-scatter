package main

import (
	"bufio"
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
	"time"
)

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

	dirPath := CFTablesMap["out_dir"].(string)
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
