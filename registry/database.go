package registry

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var schemaCollection = getCollection()

func initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// InitializeDatabase mongodb database
func getCollection() *mongo.Collection {
	initialize()
	mongoDBConnectionString := os.Getenv("MONGODB_CONNECTION_STRING")
	if mongoDBConnectionString == "" {
		log.Fatal("MONGODB_CONNECTION_STRING should not be empty")
	}

	clientOptions := options.Client().ApplyURI(mongoDBConnectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	return client.Database("schemaRegistry").Collection("schemas")
}

func saveSchema(s *SchemaRegistry) error {
	upsert := true
	updateOptions := options.UpdateOptions{Upsert: &upsert}
	update := bson.M{
		"$set": s,
	}

	filter := bson.D{{Key: "servicename", Value: s.ServiceName}}

	_, err := schemaCollection.UpdateOne(context.TODO(), filter, update, &updateOptions)
	if err != nil {
		return err
	}

	return nil
}

func getServiceSchemas(results *[]string) error {
	findOptions := options.Find()
	serviceSchemas, err := schemaCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for serviceSchemas.Next(context.TODO()) {
		var elem SchemaRegistry
		err := serviceSchemas.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		value, _ := json.Marshal(elem)
		err = redisDB.Set(ctx, elem.ServiceName, value, 0).Err()
		if err != nil {
			return err
		}

		*results = append(*results, elem.ServiceName)
	}

	if err := serviceSchemas.Err(); err != nil {
		log.Fatal(err)
	}

	serviceSchemas.Close(context.TODO())

	return nil
}
