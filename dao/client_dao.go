package dao

import (
	"HealthRobotServer/models"
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

	if !d.engine.HasTable("client_info") {
		if err := d.engine.Table("client_info").CreateTable(&models.ClientInfo{}).Error; err != nil {
			log.Println("client_dao.go Insert() Creat table error: ",err)
			return err
		}
	}

	if err := d.engine.Table("client_info").Create(&info).Error; err != nil {
		log.Println("client_dao.go Insert() Insert error: ",err)
		return err
	}

	return nil
}

func (d *ClientDao) Get(id int) *models.ClientInfo {
	return nil
}

//返回在线医生列表
func (d *ClientDao) GetAll() []models.ClientInfo {
	var clients []models.ClientInfo
	d.engine.Table("client_info").Where("online_status = ? AND client_type = ?","1", "doctor").Find(&clients)
	if clients == nil {
		log.Println("client_dao.go GetAll() no doctor is online")
	}
	return clients
}

func (d* ClientDao) Search(account string) *models.ClientInfo {
	var info = models.ClientInfo{
		Base:models.Base{
			ID:0,
		},
	}
	d.engine.Table("client_info").Where("client_account = ?", account).Find(&info)
	if info.ID == 0 {
		log.Println("client_dao.go Search() the client does not register")
		return nil
	}
	return &info
}

func (d* ClientDao) Update(info *models.ClientInfo) error {
	if err := d.engine.Table("client_info").Where("client_account = ?", info.ClientAccount).Updates(info).Error; err != nil {
		log.Println("client_dao.go Update() update client error",err)
		return err
	}
	log.Println("client_dao.go Update() update client info successfully")
	return nil
}
