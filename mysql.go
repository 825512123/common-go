package common_go

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
