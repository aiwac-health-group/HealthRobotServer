package dao

import (
	"HealthRobotServer/datasource"
	"github.com/jinzhu/gorm"
	"log"
)

type Dao struct {
	Engine *gorm.DB
}

func NewDao(Engine *gorm.DB) *Dao {
	return &Dao {
		Engine:Engine,
	}
}

//事务处理
func (d *Dao) TransactionBegin() {
	//切换数据库连接为事务连接
	d.Engine = d.Engine.Begin()
}

func (d *Dao) RollBack() {
	d.Engine.Rollback()
}

func (d *Dao) Commit() {
	d.Engine.Commit()
	d.Engine = datasource.Instance()
}

func (d *Dao) Insert(tableName string, item interface{}) error {

	if !d.Engine.HasTable(tableName) {
		if err := d.Engine.Table(tableName).CreateTable(item).Error; err != nil {
			log.Printf("dboperation.go Insert() Creat %s table error: %s",tableName, err)
			return err
		}
	}

	if err := d.Engine.Table(tableName).Create(item).Error; err != nil {
		log.Printf("dboperation.go Insert() Insert %s error: %s", item, err)
		return err
	}

	return nil
}

func (d *Dao) GetList(result interface{}, tableName string, query string, condition...interface{}) {
	d.Engine.Table(tableName).Where(query,condition...).Find(result)
	if result == nil {
		log.Println("dboperation.go GetList() no match query")
	}
}

func (d *Dao) Search(result interface{},tableName string, query string, condition...interface{}) error {
	err := d.Engine.Table(tableName).Where(query, condition...).Find(result).Error
	if err != nil {
		log.Println("dboperation.go Search() no match query: ", err)
	}
	return err
}

func (d *Dao) PreLoadSearch(result interface{}, preLoad string, tableName string, query string, condition...interface{}) error {
	err := d.Engine.Table(tableName).Where(query, condition...).Preload(preLoad).Find(result).Error
	if err != nil {
		log.Println("dboperation.go PreLoadSearch() no match query: ", err)
	}
	return err
}

func (d* Dao) Update(item interface{}, tableName string, query string, condition...interface{}) error {
	if err := d.Engine.Table(tableName).Where(query, condition...).Updates(item).Error; err != nil {
		log.Println("dboperation.go Update() update error",err)
		return err
	}
	log.Println("dboperation.go Update() update successfully")
	return nil
}