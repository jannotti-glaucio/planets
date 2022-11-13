package util

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToUpdate converts interface to bson.D with $set
func ToUpdate(v interface{}) (*primitive.D, error) {
	bsonData, err := ToDoc(v)
	if err != nil {
		return nil, err
	}

	return &bson.D{{"$set", bsonData}}, nil
}

// ToDoc converts interface to bson.D
func ToDoc(v interface{}) (*bson.D, error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	var doc *bson.D
	err = bson.Unmarshal(data, &doc)
	return doc, err
}
