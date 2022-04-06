package models

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const scoreToWin = 121

type Game struct {
    Left Player
	Right Player
	LeftScore int `bson:"leftscore"`
	RightScore int `bson:"rightscore"`
}

func (game Game) isBeingPlayed() bool {
	return game.LeftScore < scoreToWin && game.RightScore < scoreToWin
}

func NewGame(client *mongo.Client, left Player, right Player) interface{} {
	collection := client.Database("Cribbage").Collection("games")
	game := Game{Left: left, Right: right, LeftScore: 0, RightScore: 0}
	insertResult, err := collection.InsertOne(context.TODO(), game)
	if err != nil {
		log.Fatal(err)
	}
	return insertResult.InsertedID
}

func RemoveGame(client *mongo.Client, id interface{}) {
	collection := client.Database("Cribbage").Collection("games")
	filter := bson.D{{"_id", id}}
	deleteRes, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted", deleteRes.DeletedCount, "games with ID =", id)
}

func UpdateScore(client *mongo.Client, id interface{}, leftScore int, rightScore int) interface{} {
	collection := client.Database("Cribbage").Collection("games")
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"leftscore", leftScore},{"rightscore", rightScore}}}}
	updateRes, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return updateRes.UpsertedID
}

func PrintGames(client *mongo.Client) {
	collection := client.Database("Cribbage").Collection("games")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	if len(results) == 0 {
		fmt.Println("No Games Exist")
	}
	for _, result := range results {
		fmt.Println(result)
	}
	fmt.Println()
}

func (game *Game) PrintStats() {
	if game.isBeingPlayed() {
		fmt.Println(game.Left.Name, "vs", game.Right.Name, "\b! The score is", game.LeftScore, "to", game.RightScore, "\b!")
	} else {
		leftWin := game.LeftScore > game.RightScore
		if leftWin {
			fmt.Println(game.Left.Name, "beats", game.Right.Name, "\b! With a score of", game.LeftScore, "to", game.RightScore, "\b!")
		} else {
			fmt.Println(game.Right.Name, "beats", game.Left.Name, "\b! With a score of", game.RightScore, "to", game.LeftScore, "\b!")
		}
	}
}