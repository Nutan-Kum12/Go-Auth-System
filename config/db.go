package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

//global database object
//later use in repository and service layer

func ConnectDB() {
	uri := os.Getenv("MONGO_URI")
	fmt.Println("Mongo URI:", uri)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	//cancelling the request after 10 seconds if db not respond
	defer cancel()
	// connect mongo client
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		panic(err)
	}
	// ping database(to check mongo reachable or not)
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("MongoDB Connected")

	// database select
	DB = client.Database("authDB")
	fmt.Println("database created")
}
