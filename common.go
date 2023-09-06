package common_go

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

var OPEN_LOG int8
var tr *http.Transport

//JSON 普通返回不记录日志
func JSON(c *gin.Context, msg string, obj interface{}) {
	c.JSON(http.StatusOK, gin.H{"msg": msg, "data": obj})
}

//JSONMsg 信息返回并记录日志
func JSONMsg(c *gin.Context, msg string, obj interface{}) {
	LogMsg(msg, obj)
	JSON(c, msg, obj)
}

//JSONStruct 结构体返回并记录日志
func JSONStruct(c *gin.Context, msg string, obj interface{}) {
	LogStruct(msg, obj)
	JSON(c, msg, obj)
}

//LogStruct 结构体日志
func LogStruct(msg string, obj interface{}) {
	LogMap(msg, StructToMap(obj))
}

//LogMap map日志
func LogMap(msg string, data map[string]interface{}) {
	if OPEN_LOG == 0 {
		return
	}
	Logger().WithFields(logrus.Fields{
		"name": "map",
		"data": data,
	}).Info(msg)
}

//LogMsg 信息日志
func LogMsg(msg ...interface{}) {
	if OPEN_LOG == 0 {
		return
	}
	Logger().WithFields(logrus.Fields{
		"name": "msg",
	}).Info(msg)
}

//GetDate 获取当前日期
func GetDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//GetDateByLayout 获取当前日期
func GetDateByLayout(layout string) string {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return time.Now().Format(layout)
}

//MbSubstr 模拟PHP mb_substr 字符串截取函数支持utf8编码
func MbSubstr(str string, start, length int) string {
	rs := []rune(str)
	l := len(rs)
	if l > start+length {
		l = start + length
	}
	return string(rs[start:l])
}

//KSortMapSs map排序
func KSortMapSs(m map[string]string) map[string]string {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	res := make(map[string]string)
	for _, key := range keys {
		res[key] = m[key]
	}
	return res
}

//StructToMap 结构体转map
func StructToMap(obj interface{}) (m map[string]interface{}) {
	j, _ := json.Marshal(obj)
	_ = json.Unmarshal(j, &m)
	return
}

//StructToJsonStr 结构体转json字符串
func StructToJsonStr(obj interface{}) string {
	j, _ := json.Marshal(obj)
	var m map[string]interface{}
	_ = json.Unmarshal(j, &m)
	marshal, _ := json.Marshal(m)
	return string(marshal)
}

//MapToJsonStr map转json字符串
func MapToJsonStr(m map[string]interface{}) string {
	marshal, _ := json.Marshal(m)
	return string(marshal)
}

//JsonStrToMap json字符串转map
func JsonStrToMap(str string) (m map[string]interface{}) {
	_ = json.Unmarshal([]byte(str), &m)
	return
}

//InArray 根据need类型，判断need是否包含在haystack中 数组包含
func InArray(need interface{}, haystack interface{}) bool {
	switch key := need.(type) {
	case int:
		for _, item := range haystack.([]int) {
			if item == key {
				return true
			}
		}
	case string:
		for _, item := range haystack.([]string) {
			if item == key {
				return true
			}
		}
	case int64:
		for _, item := range haystack.([]int64) {
			if item == key {
				return true
			}
		}
	case float64:
		for _, item := range haystack.([]float64) {
			if item == key {
				return true
			}
		}
	default:
		return false
	}
	return false
}

//GetInterfaceToInt interface转int
func GetInterfaceToInt(t1 interface{}) int {
	var t2 int
	switch t1.(type) {
	case uint:
		t2 = int(t1.(uint))
		break
	case int8:
		t2 = int(t1.(int8))
		break
	case uint8:
		t2 = int(t1.(uint8))
		break
	case int16:
		t2 = int(t1.(int16))
		break
	case uint16:
		t2 = int(t1.(uint16))
		break
	case int32:
		t2 = int(t1.(int32))
		break
	case uint32:
		t2 = int(t1.(uint32))
		break
	case int64:
		t2 = int(t1.(int64))
		break
	case uint64:
		t2 = int(t1.(uint64))
		break
	case float32:
		t2 = int(t1.(float32))
		break
	case float64:
		t2 = int(t1.(float64))
		break
	case string:
		t2, _ = strconv.Atoi(t1.(string))
		break
	default:
		t2 = t1.(int)
		break
	}
	return t2
}

func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

func Sha256(str string) string {
	m := sha256.New()
	m.Write([]byte(str))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

//Md5 生成32位md5字串
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//GetRandString 获取随机字符串
func GetRandString(length, t int) string {
	if length < 1 {
		return ""
	}
	char := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	switch t {
	case 1: //"upper":
		char = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case 2: //"lower":
		char = "abcdefghijklmnopqrstuvwxyz"
	case 3: //"number":
		char = "0123456789"
	case 4: //"uppernumber":
		char = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	case 5: //"lowernumber":
		char = "abcdefghijklmnopqrstuvwxyz0123456789"
	case 6: //"lowerupper":
		char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	charArr := strings.Split(char, "")
	charLen := len(charArr)
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))

	rChar := ""
	for i := 1; i <= length; i++ {
		rChar = rChar + charArr[ran.Intn(charLen)]
	}
	return rChar
}

//GetRandInt 获取随机数
func GetRandInt(length int) string {
	if length < 1 {
		return ""
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	format := "%0" + fmt.Sprintf("%dv", length)
	l := 10
	if length > 1 {
		for i := 1; i < length; i++ {
			l *= 10
		}
	}
	return fmt.Sprintf(format, rnd.Intn(l))
}

//Post 发送post请求
//url 请求地址
//data 请求参数
//header 请求头 map[string]string{"Content-Type": "application/json;charset=utf-8"}
func Post(url string, data map[string]interface{}, header map[string]string) (result []byte, err error) {
	postData, _ := json.Marshal(data)
	//client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(postData))
	if err != nil {
		LogMsg(url, "接口请求失败")
		return nil, err
	}
	for s, s2 := range header {
		req.Header.Add(s, s2)
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   3 * time.Second, // 超时加在这里，是每次调用的超时
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)
	return
}

// MkPath 创建路径
// path 路径 当前路径下创建
// 返回 是否创建成功
func MkPath(path string) bool {
	if !PathExists(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			fmt.Println("path err:", err)
			return false
		}
		return true
	}
	return true
}

// MkFile 创建文件并写入内容
// path 文件路径
// info 文件内容
// 返回 是否创建成功
func MkFile(path, info string) bool {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("file err:", err)
		return false
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(info)
	if err != nil {
		fmt.Println("file write err:", err)
		return false
	}
	writer.Flush()
	return true
}

// 判定文件夹或文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//Logger 公共实例化log方法
func Logger() *logrus.Logger {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := now.Format("2006-01-02") + ".log"
	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//实例化
	logger := logrus.New()

	//设置输出
	logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logger
}

//func LoggerToFile() gin.HandlerFunc {
//	logger := Logger()
//	return func(c *gin.Context) {
//		// 开始时间
//		startTime := time.Now()
//
//		// 处理请求
//		c.Next()
//
//		// 结束时间
//		endTime := time.Now()
//
//		// 执行时间
//		latencyTime := endTime.Sub(startTime)
//
//		// 请求方式
//		reqMethod := c.Request.Method
//
//		// 请求路由
//		reqUri := c.Request.RequestURI
//
//		// 状态码
//		statusCode := c.Writer.Status()
//
//		// 请求IP
//		clientIP := c.ClientIP()
//
//		//日志格式
//		logger.Infof("| %3d | %13v | %15s | %s | %s |",
//			statusCode,
//			latencyTime,
//			clientIP,
//			reqMethod,
//			reqUri,
//		)
//	}
//}
