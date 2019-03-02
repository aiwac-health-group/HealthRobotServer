package dao

import (
	"HealthRobotServer/models"
	"github.com/jinzhu/gorm"
)

//TreatInfoDao定义了在线问诊处理列表的相关数据库操作

type TreatInfoDao struct {
	engine *gorm.DB
}

func NewTreatInfoDao(engine *gorm.DB) *TreatInfoDao {
	return &TreatInfoDao{
		engine:engine,
	}
}

//更新问诊请求的状态
func (d *TreatInfoDao) Update(TreatInfo *models.TreatInfo) {

}

//插入新的问诊请求
func (d *TreatInfoDao) Insert(TreatInfo *models.TreatInfo) {

}

//返回满足状态条件的问诊列表
func (d *TreatInfoDao) GetList(status string) []models.TreatInfo {
	return nil
}
