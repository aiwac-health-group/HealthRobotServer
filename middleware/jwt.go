package middleware

import (
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
)

func JwtHandler() *jwtmiddleware.Middleware {
	return jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{},error) {
			return []byte("HealthRobot Secret"),nil
		},
		SigningMethod:jwt.SigningMethodHS256,
	})
}

func GenToken()  *jwt.Token {
	return nil
}
