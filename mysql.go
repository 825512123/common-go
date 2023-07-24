package common_go

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

func InitDb(dsn string) *gorm.DB {
	//dsn := "root:root@tcp(127.0.0.1:3306)/ff_device?charset=utf8&parseTime=True&loc=Local"
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // 使用单数表名
	}}); err != nil {
		//GVA_DB = nil
		Logger().WithFields(logrus.Fields{
			"name": "hanyun",
		}).Error("mysql connect ping failed, err:", "Error")
	} else {
		//GVA_DB = db
		return db
	}
	return nil
}

// Column 数据字段类型
type Column struct {
	ColumnName    string `json:"column_name"`
	DataType      string `json:"data_type"`
	ColumnComment string `json:"column_comment"`
	ColumnKey     string `json:"column_key"`
	Extra         string `json:"extra"`
}

// 字符串转为大驼峰
func BigHump(str string) (data string) {
	arr := strings.Split(str, "_")
	for _, a := range arr {
		data += InitialToCapital(a)
	}
	return data
}

// InitialToCapital 首字母转大写
func InitialToCapital(str string) string {
	var InitialToCapitalStr string
	strRune := []rune(str)
	for i := 0; i < len(strRune); i++ {
		if i == 0 {
			if strRune[i] >= 97 && strRune[i] <= 122 {
				strRune[i] -= 32
				InitialToCapitalStr += string(strRune[i])
			} else {
				return str
			}
		} else {
			InitialToCapitalStr += string(strRune[i])
		}
	}
	return InitialToCapitalStr
}

//TableNameToTableStruct 数据库-数据表转结构体输出
func TableNameToTableStruct(db *gorm.DB, table string) {
	var columns []*Column
	db.Debug().Raw("select column_name, data_type, column_comment, column_key, extra from information_schema.columns where table_name = ? and table_schema =(select database()) order by ordinal_position ", table).Scan(&columns)
	TableToStruct(columns, table)
}

// 数据表转结构体
func TableToStruct(data []*Column, table string) {
	// ----- 拼接生成的struct  start--------
	structStr := fmt.Sprintf("type %s struct {\n", BigHump(table))
	for _, column := range data {
		structStr += "    " + BigHump(column.ColumnName) //InitialToCapital(column.ColumnName)
		if column.DataType == "tinyint" {
			structStr += " int "
		} else if column.DataType == "decimal" {
			structStr += " float64 "
		} else if column.DataType == "bigint" || column.DataType == "int" {
			structStr += " int64 "
		} else {
			structStr += " string "
		}
		structStr += fmt.Sprintf("`gorm:\"column:%s;comment('%s')\" json:\"%s\"` \n", column.ColumnName, column.ColumnComment, column.ColumnName)
		//if column.Extra != "auto_increment" {
		//	structStr += fmt.Sprintf("`gorm:\"comment('%s')\" json:\"%s\"` \n",
		//		column.ColumnComment, column.ColumnName)
		//} else {
		//	structStr += fmt.Sprintf("`gorm:\"not null comment('%s') INT(11)\" json:\"%s\"` \n",
		//		column.ColumnComment, column.ColumnName)
		//}
	}
	structStr += "} \n"
	fmt.Println(structStr)
}
