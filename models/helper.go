package models

import "gopkg.in/mgo.v2/bson"

func MatchById(id string) bson.M {
	return bson.M{
		"$match": bson.M{
			"_id": bson.ObjectIdHex(id),
		},
	}
}
