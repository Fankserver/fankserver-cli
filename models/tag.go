package models

import (
	"log"

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
	log.Println("count", count)
	if count == 0 {
		err = Db.C(TagCollection).Insert(Tag{
			Name: "admin",
		})
		if err != nil {
			panic(err)
		}
	}
}
