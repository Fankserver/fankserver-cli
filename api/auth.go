package api

import (
	"encoding/hex"
	"math/rand"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fank/validator"
	"github.com/fankserver/fankserver-cli/config"
	"github.com/fankserver/fankserver-cli/connection"
	"github.com/fankserver/fankserver-cli/models"
	"golang.org/x/crypto/sha3"
	iris "gopkg.in/kataras/iris.v6"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// NewAuthAPI will register routes for AuthAPI
func NewAuthAPI(router *iris.Router) {
	auth := AuthAPI{}
	router.Post("/login", auth.Login)
	router.Post("/register", auth.Register)
}

type AuthAPI struct{}

func (a *AuthAPI) Login(ctx *iris.Context) {
	db := ctx.Get("mongo").(connection.MongoDB)

	loginUser := models.User{}
	err := ctx.ReadJSON(&loginUser)
	if err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}

	usr := models.User{}
	err = db.C(models.UserCollection).Find(bson.M{"username": loginUser.Username}).One(&usr)
	if err != nil {
		if err == mgo.ErrNotFound {
			ctx.EmitError(iris.StatusNotFound)
		} else {
			ctx.EmitError(iris.StatusInternalServerError)
		}
		return
	}

	if a.validatePassword(&usr, loginUser.Password) {
		// Remove secrets
		usr.Password = ""

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  usr.ID,
			"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
			"user": &usr,
		})
		tokenString, err := token.SignedString([]byte(config.GetConfig().Jwt.Secret))
		if err != nil {
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}

		ctx.HTML(iris.StatusOK, tokenString)
	} else {
		ctx.EmitError(iris.StatusForbidden)
	}
}

func (a *AuthAPI) Register(ctx *iris.Context) {
	db := ctx.Get("mongo").(connection.MongoDB)

	registerUser := models.User{}
	err := ctx.ReadJSON(&registerUser)
	if err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}

	err = validator.Validate(registerUser)
	if err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}

	registerUser.Salt = a.generateSalt()
	registerUser.Password = a.hashPassword(registerUser.Password, registerUser.Salt)

	err = db.C(models.UserCollection).Insert(&registerUser)
	if err != nil {
		if db.IsDup(err) {
			ctx.EmitError(iris.StatusConflict)
		} else {
			ctx.EmitError(iris.StatusInternalServerError)
		}
		return
	}

	ctx.EmitError(iris.StatusCreated)
}

func (a *AuthAPI) validatePassword(usr *models.User, password string) bool {
	return a.hashPassword(password, usr.Salt) == usr.Password
}

var saltRunes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (a *AuthAPI) generateSalt() string {
	b := make([]byte, 128/8)
	for i := range b {
		b[i] = saltRunes[rand.Intn(len(saltRunes))]
	}
	h := sha3.New256()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func (a *AuthAPI) hashPassword(password string, salt string) string {
	h := sha3.New512()
	h.Write([]byte(password))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}
