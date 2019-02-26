package datasource

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"sync"
)

var (
	engine  *gorm.DB
	lock sync.Mutex
)

//数据库单例
func Instance() *gorm.DB {
	if engine != nil {
		return engine
	}
	lock.Lock()
	defer lock.Unlock()

	if engine != nil {
		return engine
	}

	db, err := gorm.Open("mysql", "root:aiwac2019@tcp(127.0.0.1:3306)/healthrobotdb?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		log.Fatal("dbhelper.Instance error",err)
		return nil
	}
	engine = db
	return engine
}

