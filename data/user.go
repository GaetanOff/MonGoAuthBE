package data

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User ... Struct for dealing with user data.
type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName string             `json:"fname" bson:"fname"`
	LastName  string             `json:"lname" bson:"lname"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	DOB       primitive.DateTime `json:"dob" bson:"dob"`
}

// String ... Get string representation of the User.
func (u User) String() string {
	return fmt.Sprintf("{ID:%s, fname:%s, lname:%s, email:%s, password:%s, dob:%s}",
		u.ID.String(), u.FirstName, u.LastName, u.Email, u.Password, u.DOB.Time().String())
}
