package admin

import (
	"github.com/fankserver/fankserver-cli/connection"
	"github.com/fankserver/fankserver-cli/models"
	iris "gopkg.in/kataras/iris.v6"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// NewUser will register routes for User
func NewUser(router *iris.Router) {
	user := User{}
	router.Get("/", user.GetUsers)
	router.Get("/:id", user.GetUser)
	router.Post("/:id/tag/:tag", user.AddTag)
	router.Delete("/:id/tag/:tag", user.RemoveTag)
	router.Put("/:id", user.SetUser)
}

type User struct{}

func (a *User) GetUsers(ctx *iris.Context) {
	db := ctx.Get("mongo").(connection.MongoDB)

	users := []models.User{}
	err := db.C(models.UserCollection).Find(nil).Select(bson.M{"password": 0}).All(&users)
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
	} else {
		ctx.JSON(iris.StatusOK, users)
	}
}

func (a *User) GetUser(ctx *iris.Context) {
	db := ctx.Get("mongo").(connection.MongoDB)

	var user models.User
	operations := []bson.M{
		models.MatchByID(ctx.Param("id")),
		models.UserLookupTags,
	}
	err := db.C(models.UserCollection).Pipe(operations).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			ctx.EmitError(iris.StatusNotFound)
		} else {
			ctx.EmitError(iris.StatusInternalServerError)
		}
		return
	}

	ctx.JSON(iris.StatusOK, user)
}

func (a *User) SetUser(ctx *iris.Context) {
	db := ctx.Get("mongo").(connection.MongoDB)

	var user models.User
	err := db.C(models.UserCollection).FindId(bson.ObjectIdHex(ctx.Param("id"))).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			ctx.EmitError(iris.StatusNotFound)
		} else {
			ctx.EmitError(iris.StatusInternalServerError)
		}
		return
	}

	oldPassword := user.Password

	err = ctx.ReadJSON(&user)
	if err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}

	newPassword := user.Password

	if oldPassword != newPassword {
		// DO SOMETHING
	}

	change := mgo.Change{
		Update:    &user,
		ReturnNew: true,
	}
	var updatedUser models.User
	_, err = db.C(models.UserCollection).FindId(bson.ObjectIdHex(ctx.Param("id"))).Apply(change, &updatedUser)
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
	} else {
		ctx.EmitError(iris.StatusOK)
	}
}

func (a *User) AddTag(ctx *iris.Context) {
	db := ctx.Get("mongo").(connection.MongoDB)

	change := bson.M{
		"$push": bson.M{
			"tags": bson.ObjectIdHex(ctx.Param("tag")),
		},
	}
	err := db.C(models.UserCollection).UpdateId(bson.ObjectIdHex(ctx.Param("id")), change)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, err)
		return
	}

	ctx.EmitError(iris.StatusOK)
}

func (a *User) RemoveTag(ctx *iris.Context) {
	db := ctx.Get("mongo").(connection.MongoDB)

	change := bson.M{
		"$pull": bson.M{
			"tags": bson.ObjectIdHex(ctx.Param("tag")),
		},
	}
	err := db.C(models.UserCollection).UpdateId(bson.ObjectIdHex(ctx.Param("id")), change)
	if err != nil {
		if err == mgo.ErrNotFound {
			ctx.EmitError(iris.StatusNotFound)
		} else {
			ctx.EmitError(iris.StatusInternalServerError)
		}
		return
	}

	ctx.EmitError(iris.StatusOK)
}
