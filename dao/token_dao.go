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

	if !d.engine.HasTable("token_info") {
		if err := d.engine.Table("token_info").CreateTable(&models.Token{}).Error; err != nil {
			log.Println("token_dao.go Insert() creat table error",err)
			return err
		}
		log.Println("token_dao.go Insert() Creat table Successfully")
	}

	if err := d.engine.Table("token_info").Create(token).Error; err != nil {
		log.Fatal("token_dao.go Insert() insert error",err)
		return err
	}
	log.Println("token_dao.go Insert() Insert token info successfully")
	return nil
}

func (d *TokenDao) Update(token *models.Token) error {
	if err := d.engine.Table("token_info").Where("client_account = ?", token.ClientAccount).Updates(token).Error; err != nil {
		log.Println("token_dao.Update update token error",err)
		return err
	}
	log.Println("token_dao.go Update() update token info successfully")
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
		log.Println("token_dao.go Search() the token does not exist yet")
		return nil
	}
	return &token
}

