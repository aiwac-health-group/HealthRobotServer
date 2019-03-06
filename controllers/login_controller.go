package controllers

import (
	"HealthRobotServer-master/constants"
	"HealthRobotServer-master/middleware"
	"HealthRobotServer-master/models"
	"HealthRobotServer-master/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"strings"
)

type LoginController struct {
	Ctx iris.Context
	Service services.LoginService
}

func (c *LoginController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle("GET","/","Welcome")
	b.Handle("POST","/","LoginWithPassword")
	//b.Handle("POST", "/loginWithIdentifyCode", "LoginWithIdentifyCode")
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

	//根据账号验证用户是否存在
	if clientInfo := c.Service.SearchClientInfo(loginInfo.Account); clientInfo.ID == 0 {
		_, _ = c.Ctx.JSON(models.LoginResponse{
			LoginFlag:"failed",
		})
		return
	} else {
		if strings.EqualFold(loginInfo.Password, clientInfo.ClientPassword) {
			//密码正确
			//获取该用户的详细信息
			profile := c.Service.SearchWebClientProfile(loginInfo.Account)

			//生成token返回至用户,并更新数据库中的token
			tokenString := middleware.GenerateToken(clientInfo.ClientAccount, clientInfo.ClientType)

			_, _ = c.Ctx.JSON(models.LoginResponse{
				LoginFlag:"success",
				ClientType:profile.ClientType,
				ClientName:profile.ClientName,
				Token:tokenString,
			})

		} else {
			_, _ = c.Ctx.JSON(models.LoginResponse{
				LoginFlag:"failed",
			})
		}
		return
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
