package g

import (
	"github.com/spf13/viper"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type DBPool struct {
	Uic    *gorm.DB
	Portal *gorm.DB
}

var (
	dbp DBPool
)

func Con() DBPool {
	return dbp
}

func InitDBPool() (err error) {
	uicd, err := gorm.Open(gmysql.Open(viper.GetString("db.uic.dsn")), &gorm.Config{})
	if err != nil {
		panic(err)
		return
	}
	uicdDB, err := uicd.DB()
	uicdDB.SetMaxIdleConns(10)
	uicdDB.SetMaxOpenConns(100)
	uicdDB.SetConnMaxLifetime(time.Minute)
	dbp.Uic = uicd

	dashd, err := gorm.Open(gmysql.Open(viper.GetString("db.portal.dsn")), &gorm.Config{})
	if err != nil {
		panic(err)
		return
	}
	dashdDB, err := dashd.DB()
	dashdDB.SetMaxIdleConns(10)
	dashdDB.SetMaxOpenConns(100)
	dashdDB.SetConnMaxLifetime(time.Minute)
	dbp.Portal = dashd

	return
}
