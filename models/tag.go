package models

import (
	"github.com/fankserver/fankserver-cli/connection"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// TagCollection is a static for the name
const TagCollection = "tags"

type Tag struct {
	ID                     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name                   string        `json:"name" bson:"name" validate:"min=5"`
	TeamspeakServerGroupID int           `json:"teamspeak_sgid" bson:"teamspeak_sgid,omitempty"`
}

func indexTag() {
	Db := connection.MongoDB{}
	Db.Init()
	defer Db.Close()

	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	}
	err := Db.C(TagCollection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	count, err := Db.C(TagCollection).Find(nil).Count()
	if err != nil {
		panic(err)
	}
	if count == 0 {
		initialTagCreation(Db)
	}
}

func initialTagCreation(db connection.MongoDB) {
	tags := []Tag{
		Tag{
			Name: "admin",
		},
	}

	for tag := range tags {
		err := db.C(TagCollection).Insert(tag)
		if err != nil {
			panic(err)
		}
	}
}