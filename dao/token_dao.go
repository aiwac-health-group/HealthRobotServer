package dao

import (
	"HealthRobotServer/models"
	"github.com/jinzhu/gorm"
	"log"
)

type TokenDao struct {
	engine *gorm.DB
}

func NewTokenDao(engine *gorm.DB) *TokenDao {
	return &TokenDao{
		engine:engine,
	}
}

func (d *TokenDao) Insert(token *models.Token) error {

	if !d.engine.HasTable(&models.Token{}) {
		if err := d.engine.Table("token_info").CreateTable(&models.Token{}).Error; err != nil {
			log.Fatal("token_dao.Insert creat table error",err)
			return err
		}
	}

	if err := d.engine.Table("token_info").Create(token).Error; err != nil {
		log.Fatal("token_dao.Insert insert error",err)
		return err
	}

	return nil
}

func (d *TokenDao) Update(token *models.Token) error {
	if err := d.engine.Table("token_info").Updates(token).Error; err != nil {
		log.Println("token_dao.Update update token error",err)
		return err
	}
	return nil
}

func (d *TokenDao) Search(raw string) *models.Token {
	var token = models.Token{
		Base:models.Base{
			ID:0,
		},
	}
	d.engine.Table("token_info").Where("raw_token = ?",raw).Find(&token)
	if token.ID == 0 {
		log.Println("admin_dao.Search the client does not register")
		return nil
	}
	return &token
}

