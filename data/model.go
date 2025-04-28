package data

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	UserId   bson.ObjectID `bson:"_id,omitempty"`
	Username string        `bson:"username"`
	Email    string        `bson:"email"`
	Password string        `bson:"password"`
}

type Note struct {
	NoteId    bson.ObjectID `bson:"_id,omitempty"`
	UserId    string        `bson:"userId"`
	Title     string        `bson:"title"`
	Text      string        `bson:"text"`
	CreatedAt string        `bson:"createdAt"`
	UpdatedAt string        `bson:"updatedAt"`
}
