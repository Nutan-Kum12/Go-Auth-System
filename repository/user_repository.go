package repository

import (
	"Auth/config"
	"Auth/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	SaveUser(user model.User) error
	FindUserByEmail(email string) (model.User, error)
}
type MongoRepository struct{}

func (m MongoRepository) SaveUser(user model.User) error {
	collection := config.DB.Collection("users")
	context, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()
	_, err := collection.InsertOne(
		context,
		user,
	)
	return err
}
func (m MongoRepository) FindUserByEmail(email string) (model.User, error) {
	var user model.User
	collection := config.DB.Collection("users")
	context, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()
	err := collection.FindOne(
		context,
		bson.M{"email": email}, //Mongodb query object(email match)
	).Decode(&user) //convert bson to go user object
	return user, err
}
