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
	Db := connection.MongoDB{}
	Db.Init()
	defer Db.Close()

	users := []models.User{}
	err := Db.C(models.UserCollection).Find(nil).Select(bson.M{"password": 0}).All(&users)
	if err != nil {
		ctx.HTML(iris.StatusInternalServerError, "Error")
	} else {
		ctx.JSON(iris.StatusOK, users)
	}
}

func (a *User) GetUser(ctx *iris.Context) {
	Db := connection.MongoDB{}
	Db.Init()
	defer Db.Close()

	var user models.User
	operations := []bson.M{
		models.MatchById(ctx.Param("id")),
		models.UserLookupTags,
	}
	err := Db.C(models.UserCollection).Pipe(operations).One(&user)
	if err != nil {
		ctx.HTML(iris.StatusNotFound, "Not found")
	} else {
		ctx.JSON(iris.StatusOK, user)
	}
}

func (a *User) SetUser(ctx *iris.Context) {
	Db := connection.MongoDB{}
	Db.Init()
	defer Db.Close()

	var user models.User
	err := Db.C(models.UserCollection).FindId(bson.ObjectIdHex(ctx.Param("id"))).One(&user)
	if err != nil {
		ctx.HTML(iris.StatusNotFound, "Not found")
		return
	}

	oldPassword := user.Password

	err = ctx.ReadJSON(&user)
	if err != nil {
		ctx.HTML(iris.StatusBadRequest, "error parsing")
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
	_, err = Db.C(models.UserCollection).FindId(bson.ObjectIdHex(ctx.Param("id"))).Apply(change, &updatedUser)
	if err != nil {
		ctx.HTML(iris.StatusInternalServerError, "")
	} else {
		ctx.EmitError(iris.StatusOK)
	}
}

func (a *User) AddTag(ctx *iris.Context) {
	Db := connection.MongoDB{}
	Db.Init()
	defer Db.Close()

	change := bson.M{
		"$push": bson.M{
			"tags": bson.ObjectIdHex(ctx.Param("tag")),
		},
	}
	err := Db.C(models.UserCollection).UpdateId(bson.ObjectIdHex(ctx.Param("id")), change)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, err)
		return
	}

	ctx.EmitError(iris.StatusOK)
}

func (a *User) RemoveTag(ctx *iris.Context) {
	Db := connection.MongoDB{}
	Db.Init()
	defer Db.Close()

	change := bson.M{
		"$pull": bson.M{
			"tags": bson.ObjectIdHex(ctx.Param("tag")),
		},
	}
	err := Db.C(models.UserCollection).UpdateId(bson.ObjectIdHex(ctx.Param("id")), change)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, err)
		return
	}

	ctx.EmitError(iris.StatusOK)
}
