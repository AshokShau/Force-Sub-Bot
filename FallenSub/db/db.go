package db

import (
	"context"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	tdCtx = context.TODO()

	fSubCell *mongo.Collection
)

func init() {
	mongoClient, err := mongo.Connect(tdCtx, options.Client().ApplyURI(config.DatabaseURI))
	if err != nil {
		log.Fatalf("[Database][Connect]: %v", err)
	}

	fSubCell = mongoClient.Database(config.DbName).Collection("forceSub")

}

// updateOne func to update one document
func updateOne(collection *mongo.Collection, filter bson.M, data interface{}) (err error) {
	_, err = collection.UpdateOne(tdCtx, filter, bson.M{"$set": data}, options.Update().SetUpsert(true))
	if err != nil {
		config.ErrorLog.Printf("[Database][updateOne]: %v", err)
	}
	return
}

// findOne func to find one document
func findOne(collection *mongo.Collection, filter bson.M) (res *mongo.SingleResult) {
	return collection.FindOne(tdCtx, filter)
}
