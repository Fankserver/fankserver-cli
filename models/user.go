package models

import (
	"github.com/fankserver/fankserver-cli/connection"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserCollection is a static for the name
const UserCollection = "users"

// UserLookupTags predefined lookup for tags in user
var UserLookupTags = bson.M{
	"$lookup": bson.M{
		"from":         TagCollection,
		"localField":   "tags",
		"foreignField": "_id",
		"as":           "tags",
	},
}

type User struct {
	ID               bson.ObjectId     `json:"id" bson:"_id,omitempty"`
	Username         string            `json:"username" bson:"username" validate:"regexp=^[\\w\\-]{4\\,}$"`
	Password         string            `json:"password,omitempty" bson:"password" validate:"regexp=^\\S{6\\,}$"`
	Salt             string            `json:"-" bson:"salt"`
	EMail            string            `json:"email" bson:"email" validate:"regexp=^(([^<>()[\\]\\\\.\\,;:\\s@\"]+(\\.[^<>()[\\]\\\\.\\,;:\\s@\"]+)*)|(\".+\"))@((\\[[0-9]{1\\,3}\\.[0-9]{1\\,3}\\.[0-9]{1\\,3}\\.[0-9]{1\\,3}])|(([a-zA-Z\\-0-9]+\\.)+[a-zA-Z]{2\\,}))$"`
	Avatar           string            `json:"avatar" bson:"avatar"`
	ApplicationLinks []ApplicationLink `json:"appLinks" bson:",omitempty"`
	Tags             []Tag             `json:"tags" bson:"tags,omitempty"`
}

func indexUser() {
	db := connection.MongoDB{}
	db.Init()
	defer db.Close()

	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	}
	err := db.C(UserCollection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	index = mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	}
	err = db.C(UserCollection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

type ApplicationLink struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
}
