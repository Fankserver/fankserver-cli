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
	db := ctx.Get("mongo").(connection.MongoDB)

	tags := []models.Tag{}
	err := db.C(models.TagCollection).Find(nil).All(&tags)
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(iris.StatusOK, tags)
}

func (a *Tag) GetTag(ctx *iris.Context) {
	db := ctx.Get("mongo").(connection.MongoDB)

	var tag models.Tag
	err := db.C(models.TagCollection).FindId(bson.ObjectIdHex(ctx.Param("id"))).One(&tag)
	if err != nil {
		if err == mgo.ErrNotFound {
			ctx.EmitError(iris.StatusNotFound)
		} else {
			ctx.EmitError(iris.StatusInternalServerError)
		}
		return
	}

	ctx.JSON(iris.StatusOK, tag)
}

func (a *Tag) AddTag(ctx *iris.Context) {
	db := ctx.Get("mongo").(connection.MongoDB)

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

	err = db.C(models.TagCollection).Insert(tag)
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

func (a *Tag) SetTag(ctx *iris.Context) {
	db := ctx.Get("mongo").(connection.MongoDB)
	objectID := bson.ObjectIdHex(ctx.Param("id"))

	var tag models.Tag
	err := db.C(models.TagCollection).FindId(objectID).One(&tag)
	if err != nil {
		if err == mgo.ErrNotFound {
			ctx.EmitError(iris.StatusNotFound)
		} else {
			ctx.EmitError(iris.StatusInternalServerError)
		}
		return
	}

	err = ctx.ReadJSON(&tag)
	if err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}

	change := mgo.Change{
		Update:    &tag,
		ReturnNew: true,
	}
	var updatedTag models.Tag
	_, err = db.C(models.TagCollection).FindId(objectID).Apply(change, &updatedTag)
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	ctx.EmitError(iris.StatusOK)
}

func (a *Tag) RemoveTag(ctx *iris.Context) {
	db := ctx.Get("mongo").(connection.MongoDB)

	err := db.C(models.TagCollection).RemoveId(bson.ObjectIdHex(ctx.Param("id")))
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
