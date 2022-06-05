package common_go

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"time"
)

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
	Logger().WithFields(logrus.Fields{
		"name": "map",
		"data": data,
	}).Info(msg)
}

//LogMsg 信息日志
func LogMsg(msg ...interface{}) {
	Logger().WithFields(logrus.Fields{
		"name": "msg",
	}).Info(msg)
}

//GetDate 获取当前日期
func GetDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
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

func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

func Post(url string, data map[string]interface{}, header map[string]string) (result []byte, err error) {
	postData, _ := json.Marshal(data)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(postData))
	if err != nil {
		LogMsg(url, "接口请求失败")
		return nil, err
	}
	for s, s2 := range header {
		req.Header.Add(s, s2)
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)
	return
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
