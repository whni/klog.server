package main

import (
	"context"
	"fmt"
	"log"
	//"time"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoDBHosts = "localhost:27017"
	AuthDatabase = "klog"
	AuthUserName = "klog_user"
	AuthPassword = "klog_pwd"
)

func main() {
	// Set client options
	dbURL := fmt.Sprintf("mongodb://%s:%s@%s/%s", AuthUserName, AuthPassword, MongoDBHosts, AuthDatabase)
	clientOptions := options.Client().ApplyURI(dbURL)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	defer func() {
		err = client.Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection to MongoDB closed.")
	}()

	db := client.Database(AuthDatabase)
	collectNames, err := db.ListCollectionNames(context.TODO(), nil)
	log.Println("collection names: ", collectNames)
	if err != nil {
		log.Fatal(err)
	}
	collectNames = []string{"institutes", "teachers"}

	for _, cName := range collectNames {
		collections := db.Collection(cName)
		cCount, err := collections.EstimatedDocumentCount(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Collection: %s: count %d", cName, cCount)
	}

}
