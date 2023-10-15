package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type AuthorWithoutId struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func main() {
	connectionURI := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("moviebox").Collection("author")
	authorID, err := addNewAuthor(coll, &AuthorWithoutId{Name: "Fedya", Age: 18})
	if err != nil {
		log.Fatal(err)
	}

	id, _ := primitive.ObjectIDFromHex("652c58664f592a3b5251fe4d")
	fmt.Printf("Author with name %s created", authorID.InsertedID)
	_, err = updateAuthorById(coll, id, bson.D{{"$set", bson.D{{"age", 19}}}})
	if err != nil {
		log.Fatal(err)
	}
}

func addNewAuthor(collection *mongo.Collection, author *AuthorWithoutId) (*mongo.InsertOneResult, error) {
	authorID, err := collection.InsertOne(context.TODO(), *author)
	if err != nil {
		return nil, err
	}
	return authorID, nil
}

func updateAuthorById(collection *mongo.Collection, id primitive.ObjectID, update bson.D) (*mongo.UpdateResult, error) {
	result, err := collection.UpdateOne(context.TODO(), bson.D{{"_id", id}}, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
