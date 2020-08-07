package utils

import "go.mongodb.org/mongo-driver/bson"

func SingleBson(k string, val string) interface{} {
	return bson.M{k: val}
}
