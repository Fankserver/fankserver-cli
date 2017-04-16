package api

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fankserver/fankserver-cli/config"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	iris "gopkg.in/kataras/iris.v6"
)

var jwtSigningMethod = jwt.SigningMethodHS256

func jwtMiddleware() iris.Handler {
	return jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetConfig().Jwt.Secret), nil
		},
		SigningMethod: jwtSigningMethod,
	})
}
