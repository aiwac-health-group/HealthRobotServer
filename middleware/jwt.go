package middleware

import (
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"time"
)

func JwtHandler() *jwtmiddleware.Middleware {
	return jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{},error) {
			return []byte("HealthRobot Secret"),nil
		},
		SigningMethod:jwt.SigningMethodHS256,
	})
}

func JwtHandlerWS() *jwtmiddleware.Middleware {
	return jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{},error) {
			return []byte("HealthRobot Secret"),nil
		},
		SigningMethod:jwt.SigningMethodHS256,
		Extractor: func(ctx iris.Context) (s string, e error) { //自定义token提取函数:从url参数中提取
			tokenValue := ctx.FormValue("token")
			if tokenValue == "" {
				return "", errors.New("token does not exist")
			}
			return tokenValue, nil
		},
	})
}

func GenerateToken(account string, clientType string) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"Account":account,
		"ClientType":clientType,
		"ExpressIn":time.Now().AddDate(0,0,1).Unix(),
	})
	tokenString, _ := jwtToken.SignedString([]byte("HealthRobot Secret"))
	return tokenString
}