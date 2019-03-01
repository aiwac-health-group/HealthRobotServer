package dao

import (
	"HealthRobotServer/models"
	"github.com/jinzhu/gorm"
)

//TreatDao定义了在线问诊处理列表的相关数据库操作

type TreatDao struct {
	engine *gorm.DB
}

func NewTreatMessageDao(engine *gorm.DB) *TreatDao {
	return &TreatDao{
		engine:engine,
	}
}

//删除处理完成的问诊请求
func (d *TreatDao) Delete(treat *models.Treat) {

}

//插入新的问诊请求
func (d *TreatDao) Insert(treat *models.Treat) {

}

//返回未处理的请求
func (d *TreatDao) GetList() []models.Treat {
	return nil
}
