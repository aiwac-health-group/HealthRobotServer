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
			ErrorCode:"401",
			ErrorDesc:"Unauthorized",
		})
		ctx.StopExecution()
		log.Fatal("JWT invalidated",err)
		return
	}
	// If everything ok then call next.
	ctx.Next()
}

func (m *AuthToken) CheckJWT(ctx iris.Context) error {
	jwtToken := ctx.Values().Get("jwt").(*jwt.Token)
	token := m.Service.GetToken(jwtToken.Raw)
	if token != nil && token.ExpressIn < time.Now().Unix() {
		//后期还可以根据请求url来判断用户是否有访问该页面的权限，比如只有管理员用户可以访问/admin路由下的资源，避免错误
		return nil
	}
	return errors.New("Unauthorized")
}