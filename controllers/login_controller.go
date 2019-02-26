package controllers

import (
	"HealthRobotServer/models"
	"HealthRobotServer/services"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"strings"
)

type LoginController struct {
	Ctx iris.Context
	Service services.ClientService
}

func (c *LoginController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle("GET","/loginWithPassword","LoginWithPassword")
}

//处理Web端用户通过工号和密码登录系统
//获得json字段中的clientID，根据clientID查询admin_info、doctor_info和service_info表
//如果查找到对应用户，则验证密码是否正确：正确，则生成token、clientType字段，返回至用户；并更新对应数据库表
//否则，返回错误至用户
func (c *LoginController) LoginWithPassword() {

	var loginInfo models.LoginRequest

	if err := c.Ctx.ReadJSON(&loginInfo); err != nil {
		fmt.Println("fail to encode request")
		return
	}

	//校验参数是否有误
	if len(loginInfo.ClientID) == 0 || len(loginInfo.Password) == 0 {
		_, _ = c.Ctx.JSON(models.LoginResponse{
			BaseResponse:models.BaseResponse{
				ErrorCode:"400",
				ErrorDesc:"wrong client id or password input",
			},
		})
		return
	}

	//验证用户是否存在
	if client := c.Service.GetByName(loginInfo.ClientID); client == nil {
		_, _ = c.Ctx.JSON(models.LoginResponse{
			BaseResponse:models.BaseResponse{
				ErrorCode:"400",
				ErrorDesc:"wrong client id or password input",
			},
		})
		return
	} else {
		if strings.EqualFold(loginInfo.Password, client.ClientPassword) {
			//密码正确,生成token返回至用户,并更新数据库中的token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
				"clientID":loginInfo.ClientID,
			})
			tokenString, _ := token.SignedString([]byte("Aiwac Secert"))
			_, _ = c.Ctx.JSON(models.LoginResponse{
				BaseResponse:models.BaseResponse{
					ErrorCode:"200",
					ErrorDesc:"login successfully",
				},
				ClientType:client.ClientType,
				Token:tokenString,
			})
		} else {
			_, _ = c.Ctx.JSON(models.LoginResponse{
				BaseResponse:models.BaseResponse{
					ErrorCode:"400",
					ErrorDesc:"wrong client id or password input",
				},
			})
		}
		return
	}

}