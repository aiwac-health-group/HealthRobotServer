package dao

import (
	"HealthRobot/models"
	"github.com/jinzhu/gorm"
	"log"
)

type ClientDao struct {
	engine *gorm.DB
}

func NewClientDao(engine *gorm.DB) *ClientDao {
	return &ClientDao {
		engine:engine,
	}
}

func (d *ClientDao) Insert(info *models.ClientInfo) error {

	if !d.engine.HasTable(&models.ClientInfo{}) {
		if err := d.engine.Table("admin_info").CreateTable(&models.ClientInfo{}).Error; err != nil {
			log.Fatal("admin_dao.Insert creat table error",err)
			return err
		}
	}

	if err := d.engine.Create(info).Error; err != nil {
		log.Fatal("admin_dao.Insert insert error",err)
		return err
	}

	return nil
}

func (d *ClientDao) Get(id int) *models.ClientInfo {
	return nil
}

func (d *ClientDao) GetAll(id int) []models.ClientInfo {
	return nil
}

func (d* ClientDao) Search(name string) *models.ClientInfo {
	var info = models.ClientInfo{
		Base:models.Base{
			ID:0,
		},
	}
	d.engine.Table("admin_info").Where("admin_name = ?",name).Find(&info)
	if info.ID == 0 {
		log.Println("admin_dao.Search the client does not register")
		return nil
	}
	return &info
}

func (d* ClientDao) Update(data *models.ClientInfo) error {
	return nil
}
