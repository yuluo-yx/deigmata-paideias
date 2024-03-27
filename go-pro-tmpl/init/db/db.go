package db

import (
	"fmt"

	"api-registry/pkg/log"
	"api-registry/pkg/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	config := model.GetConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.System.Mysql.User,
		config.System.Mysql.Password,
		config.System.Mysql.Host,
		config.System.Mysql.Port,
		config.System.Mysql.Databases,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.SysLog.Warn("Init DB Failed!")
		panic(fmt.Errorf("failed to connect to database: %s", err))
	}
	log.SysLog.Info("Init DB Success!")

}

func GetDB() *gorm.DB {

	return DB
}
