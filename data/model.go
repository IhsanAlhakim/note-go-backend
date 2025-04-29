package data

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
if using mongo-driver v2, use bson.ObjectID for document _id
if using mongo-driver, use primitive.ObjectID
*/

type User struct {
	UserId   primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type Note struct {
	NoteId    primitive.ObjectID `bson:"_id,omitempty"`
	UserId    string             `bson:"userId"`
	Title     string             `bson:"title"`
	Text      string             `bson:"text"`
	CreatedAt string             `bson:"createdAt"`
	UpdatedAt string             `bson:"updatedAt"`
}
