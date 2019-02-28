package middleware

import (
	"HealthRobotServer/models"
	"HealthRobotServer/services"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"log"
	"time"
)

type AuthToken struct {
	Service services.TokenService
	Config config
}

type config struct {

}

func NewAuthToken () *AuthToken {
	return &AuthToken{
		Service:services.NewTokenService(),
	}
}

func (m *AuthToken) Serve(ctx iris.Context) {
	if err := m.CheckJWT(ctx); err != nil {
		_, _ = ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"Unauthorized",
		})
		ctx.StopExecution()
		log.Println("JWT ",err)
		return
	}
	// If everything ok then call next.
	ctx.Next()
}

func (m *AuthToken) CheckJWT(ctx iris.Context) error {
	var jwtToken *jwt.Token
	if value := ctx.Values().Get("jwt"); value != nil {
		jwtToken = value.(*jwt.Token)
	}

	var token *models.Token
	if jwtToken == nil { //从GET请求参数中获取token串（针对websocket连接情况)
		tokenString := ctx.FormValue("token")
		token = m.Service.GetToken(tokenString)
	} else { //从HTTP请求头获取（针对http请求情况）
		token = m.Service.GetToken(jwtToken.Raw)
	}
	if token != nil && token.ExpressIn > time.Now().Unix() {
		log.Println("authorized client access")
		//后期还可以根据请求url来判断用户是否有访问该页面的权限，比如只有管理员用户可以访问/admin路由下的资源，避免错误
		return nil
	}
	return errors.New("unauthorized")
}