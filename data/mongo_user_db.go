// mongo_user_db.go

package data

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// MongoUserDB ... Object for interacting with mongodb server.
type MongoUserDB struct {
	client *mongo.Client // The driver client object
}

// NewMongoUserDB ... Create a new mongo proxy and connect to the server...
func NewMongoUserDB(ctx context.Context, uri string) (MongoUserDB, error) {
	// Setup the mongo options and connect to the server.
	clientOptions := options.Client().ApplyURI(uri)
	mongoClient, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return MongoUserDB{}, err
	}

	// Make sure we can actually talk to the server.
	err = mongoClient.Ping(ctx, nil)

	if err != nil {
		return MongoUserDB{}, err
	}

	// Okay we are good return the object the system will use to interact with the db.
	return MongoUserDB{
		mongoClient,
	}, nil
}

// GetUser ... Gets the information of a user if it exists.
func (mp *MongoUserDB) GetUser(ctx context.Context, email string) (User, error) {
	// This object will store the user data if we find it.
	user := User{}

	// Query the database...
	result := mp.client.Database("userdb").Collection("users").FindOne(ctx, bson.D{{Key:"email", Value:email}})

	// If we successfully got a user with the email...
	if result.Err() != nil {
		return user, result.Err()
	}

	// Decode the bson into a usable struct.
	result.Decode(&user)

	return user, nil
}

// CreateUser ... Create a new user in the database from the given user info.
func (mp *MongoUserDB) CreateUser(ctx context.Context, r *CreateUserRequest) (primitive.ObjectID, error) {
	// Hash the users password.
	hash, err := bcrypt.GenerateFromPassword([]byte(r.Auth.Password), 10)

	if err != nil {
		return primitive.ObjectID{}, err
	}

	// dob := time.Date(int(info.DOB.Year),
	// 	time.Month(info.DOB.Month),
	// 	int(info.DOB.Day),
	// 	0, 0, 0, 0, time.UTC)

	// Create the user object from the info.
	user := User{
		ID:        primitive.NewObjectID(),
		FirstName: r.FName,
		LastName:  r.LName,
		Email:     r.Auth.Email,
		Password:  string(hash),
		//DOB:       primitive.NewDateTimeFromTime(dob),
	}

	result, err := mp.client.Database("userdb").Collection("users").InsertOne(ctx, user)

	// Get the ID.
	if err == nil {
		log.Println("New User Created: ", result.InsertedID)
	}

	return user.ID, err
}
