package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type UserProfile struct {
	Username      string `bson:"username"`
	RaceCompleted int    `bson:"races"`
}

func createMongoClient() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
}

func DisconnectMongoClient() {
	fmt.Println("mongo client disconnected")
	client.Disconnect(context.Background())
}

func createUser(username string) {

	db := client.Database("test")
	cmd := bson.D{{"createUser", username},
		{"pwd", "pass"},
		{"roles", []bson.M{}}}

	r := db.RunCommand(context.Background(), cmd)

	if r.Err() != nil {
		fmt.Println(r.Err())
	} else {
		fmt.Println("User created successfully!")
	}

}

func insertUserProfile(name string) {
	userProfile := UserProfile{
		Username:      name,
		RaceCompleted: 0,
	}

	db := client.Database("test")
	userProfileCollection := db.Collection("user_profile")

	// Check if the username already exists
	existingUserFilter := bson.M{"username": userProfile.Username}
	existingUserCount, err := userProfileCollection.CountDocuments(context.Background(), existingUserFilter)
	if err != nil {
		fmt.Println("Error checking existing username:", err)
		return
	}

	if existingUserCount == 0 {
		_, err := userProfileCollection.InsertOne(context.Background(), userProfile)
		if err != nil {
			fmt.Println("Error inserting user profile:", err)
			return
		} else {
			fmt.Println("User profile information inserted successfully!")
		}
	} else {
		fmt.Println("user profile already exists")
	}

}

func updateRaceCompleted(user string) {
	db := client.Database("test")
	userProfileCollection := db.Collection("user_profile")

	filter := bson.M{"username": user}

	// Retrieve the existing document
	var existingUserProfile UserProfile
	var err error
	err = userProfileCollection.FindOne(context.Background(), filter).Decode(&existingUserProfile)
	if err != nil {
		fmt.Println("Error retrieving existing document:", err)
		return
	}

	// Construct the update using the existing value
	update := bson.D{
		{"$set", bson.D{
			{"races", existingUserProfile.RaceCompleted + 1},
		}},
	}

	// Update a single document
	_, err = userProfileCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("Error updating document:", err)
		return
	}

	fmt.Println("race completed")
}

func getRaceCompleted(user string) int {
	db := client.Database("test")
	userProfileCollection := db.Collection("user_profile")

	filter := bson.M{"username": user}

	// Retrieve the existing document
	var existingUserProfile UserProfile
	err := userProfileCollection.FindOne(context.Background(), filter).Decode(&existingUserProfile)
	if err != nil {
		fmt.Println("Error retrieving existing document:", err)
		return 0
	} else {
		return existingUserProfile.RaceCompleted
	}
}
