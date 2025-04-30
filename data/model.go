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
	NoteId    primitive.ObjectID `bson:"_id,omitempty" json:"noteId"`
	UserId    string             `bson:"userId" json:"userId"`
	Title     string             `bson:"title" json:"title"`
	Text      string             `bson:"text" json:"text"`
	CreatedAt string             `bson:"createdAt" json:"createdAt"`
	UpdatedAt string             `bson:"updatedAt" json:"updatedAt"`
}
