package admin

import (
	"github.com/fank/validator"
	"github.com/fankserver/fankserver-cli/connection"
	"github.com/fankserver/fankserver-cli/models"
	iris "gopkg.in/kataras/iris.v6"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// NewTag will register routes for Tag
func NewTag(router *iris.Router) {
	tag := Tag{}
	router.Get("/", tag.GetTags)
	router.Get("/:id", tag.GetTag)
	router.Post("/", tag.AddTag)
	router.Put("/:id", tag.SetTag)
	router.Delete("/:id", tag.RemoveTag)
}

type Tag struct {
	*iris.Context
}

func (a *Tag) GetTags(ctx *iris.Context) {
	Db := connection.MongoDB{}
	Db.Init()
	defer Db.Close()

	tags := []models.Tag{}
	err := Db.C(models.TagCollection).Find(nil).All(&tags)
	if err != nil {
		ctx.HTML(iris.StatusInternalServerError, "Error")
	} else {
		ctx.JSON(iris.StatusOK, tags)
	}
}

func (a *Tag) GetTag(ctx *iris.Context) {
	Db := connection.MongoDB{}
	Db.Init()
	defer Db.Close()

	var tag models.Tag
	err := Db.C(models.TagCollection).FindId(bson.ObjectIdHex(ctx.Param("id"))).One(&tag)
	if err != nil {
		ctx.HTML(iris.StatusNotFound, "Not found")
	} else {
		ctx.JSON(iris.StatusOK, tag)
	}
}

func (a *Tag) AddTag(ctx *iris.Context) {
	var tag models.Tag
	err := ctx.ReadJSON(&tag)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, err)
		return
	}

	err = validator.Validate(tag)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, err)
		return
	}

	Db := connection.MongoDB{}
	Db.Init()
	defer Db.Close()

	err = Db.C(models.TagCollection).Insert(tag)
	if err != nil {
		ctx.EmitError(iris.StatusConflict)
		return
	}

	ctx.EmitError(iris.StatusOK)
}

func (a *Tag) SetTag(ctx *iris.Context) {
	Db := connection.MongoDB{}
	Db.Init()
	defer Db.Close()

	var tag models.Tag
	err := Db.C(models.TagCollection).FindId(bson.ObjectIdHex(ctx.Param("id"))).One(&tag)
	if err != nil {
		ctx.HTML(iris.StatusNotFound, "Not found")
		return
	}

	err = ctx.ReadJSON(&tag)
	if err != nil {
		ctx.HTML(iris.StatusBadRequest, "error parsing")
		return
	}

	change := mgo.Change{
		Update:    &tag,
		ReturnNew: true,
	}
	var updatedTag models.Tag
	_, err = Db.C(models.TagCollection).FindId(bson.ObjectIdHex(ctx.Param("id"))).Apply(change, &updatedTag)
	if err != nil {
		ctx.HTML(iris.StatusInternalServerError, "")
		return
	}

	ctx.EmitError(iris.StatusOK)
}

func (a *Tag) RemoveTag(ctx *iris.Context) {
	Db := connection.MongoDB{}
	Db.Init()
	defer Db.Close()

	err := Db.C(models.TagCollection).RemoveId(bson.ObjectIdHex(ctx.Param("id")))
	if err != nil {
		ctx.HTML(iris.StatusNotFound, "Not found")
	} else {
		ctx.EmitError(iris.StatusOK)
	}
}
