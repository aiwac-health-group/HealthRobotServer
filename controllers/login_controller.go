package controllers

import (
	"HealthRobotServer/constants"
	"HealthRobotServer/middleware"
	"HealthRobotServer/models"
	"HealthRobotServer/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"strconv"
	"strings"
)

type LoginController struct {
	Ctx iris.Context
	Service services.LoginService
}

var IdentifyCode string

func (c *LoginController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle("GET","/","Welcome")
	b.Handle("POST","/","LoginWithPassword")
	b.Handle("POST", "/getIdentifyCode", "GetIdentifyCode")
	b.Handle("POST", "/loginWithIdentifyCode", "LoginWithIdentifyCode")
	b.Handle("GET","/getTokenString","GetAccessToken")
}

func (c *LoginController) Welcome() mvc.Result {
	log.Println("welcome to health Robot")
	return mvc.View{
		Name:"index.html",
	}
}

//处理Web端用户通过工号和密码登录系统
//获得json字段中的account，根据account查询client_info表
//如果查找到对应用户，则验证密码是否正确：正确，则生成token、clientType字段，返回至用户；并更新对应数据库表
//否则，返回错误至用户
func (c *LoginController) LoginWithPassword() {

	var loginInfo models.LoginRequest

	if err := c.Ctx.ReadJSON(&loginInfo); err != nil {
		log.Println("fail to encode request")
		return
	}

	log.Println("Client is logging: ", loginInfo)

	//校验参数是否有误
	if len(loginInfo.Account) == 0 || len(loginInfo.Password) == 0 {
		_, _ = c.Ctx.JSON(models.LoginResponse{
			LoginFlag:"failed",
		})
		return
	}

	var account string
	var password string
	var name string
	var clientType string

	//根据账号查找客服表及医生表来验证用户是否存在
	serviceInfo := c.Service.SearchServiceClientInfo(loginInfo.Account)
	doctorInfo := c.Service.SearchDoctorClientInfo(loginInfo.Account)
	if serviceInfo.ID == 0 && doctorInfo.ID == 0 { //账户不存在
		_, _ = c.Ctx.JSON(models.LoginResponse{
			LoginFlag:"failed",
		})
		return
	} else if serviceInfo.ID != 0 {
		account = loginInfo.Account
		password = serviceInfo.ClientPassword
		name = serviceInfo.ClientName
		clientType = serviceInfo.ClientType
	} else {
		account = loginInfo.Account
		password = doctorInfo.ClientPassword
		name = doctorInfo.ClientName
		clientType = doctorInfo.ClientType
	}

	if strings.EqualFold(loginInfo.Password, password) { //密码正确

		//生成token返回至用户,并更新数据库中的token
		tokenString := middleware.GenerateToken(account, clientType)

		_, _ = c.Ctx.JSON(models.LoginResponse{
			LoginFlag:"success",
			ClientType:clientType,
			ClientName:name,
			Token:tokenString,
		})

	} else {
		_, _ = c.Ctx.JSON(models.LoginResponse{
			LoginFlag:"failed",
		})
	}
	return

}

//机器人获取验证码
func (c *LoginController) GetIdentifyCode() {
	var request models.LoginRequest
	if err := c.Ctx.ReadJSON(&request); err != nil {
		log.Println("fail to encode request")
		return
	}

	IdentifyCode = strconv.Itoa(c.Service.GenerateIdentifyCode())
	_ = c.Service.SendIdentifyCodeToPhone(request.Account, IdentifyCode)
}

//机器人根据验证码登录
func (c *LoginController) LoginWithIdentifyCode() {
	var request models.LoginRequest
	if err := c.Ctx.ReadJSON(&request); err != nil {
		log.Println("fail to encode request")
		return
	}

	if strings.EqualFold(request.IdentifyCode, IdentifyCode) { //验证成功
		//生成token返回至用户,并更新数据库中的token
		tokenString := middleware.GenerateToken(request.Account, "robot")
		_, _ = c.Ctx.JSON(models.TokenResponse{
			BaseResponse:models.BaseResponse{
				Status:"2000",
				Message:"Login Successful",
			},
			Token:tokenString,
		})
	} else {
		_, _ = c.Ctx.JSON(models.TokenResponse{
			BaseResponse:models.BaseResponse{
				Status:"2001",
				Message:"Wrong IdentifyCode",
			},
			Token:"",
		})
	}
}

//用于机器人APP获取Token以免登陆
func (c *LoginController) GetAccessToken() {
	var request models.TokenGetRequest

	if err := c.Ctx.ReadJSON(&request); err != nil {
		log.Println("fail to encode request")
		return
	}

	if err := middleware.JwtHandler().CheckJWT(c.Ctx); err != nil {
		_, _ = c.Ctx.JSON(models.TokenResponse{
			BaseResponse:models.BaseResponse{
				Status:"2001",
				Message:"非法用户",
			},
		})
	}

	tokenString := middleware.GenerateToken(request.Account, constants.ClientType_robot)

	_, _ = c.Ctx.JSON(models.TokenResponse{
		BaseResponse:models.BaseResponse{
			Status:"2001",
			Message:"非法用户",
		},
		Token:tokenString,
	})

}
