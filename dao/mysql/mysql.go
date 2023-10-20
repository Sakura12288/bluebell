package mysql

import (
	"bluebellproject/setting"
	"fmt"

	"go.uber.org/zap"

	"gorm.io/driver/mysql"

	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(mysqlConfig *setting.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlConfig.User,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("gorm.Open failed", zap.Error(err))
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("db.DB() failed", zap.Error(err))
		return err
	}
	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConn)
	sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConn)
	return
}
func Close() {
	sqlDB, _ := db.DB()
	sqlDB.Close()
}
