package registry

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var schemaCollection *mongo.Collection

// Initialize mongodb database
func Initialize() {
	mongoDbConnectionString := os.Getenv("MONGODB_CONNECTION_STRING")
	if mongoDbConnectionString == "" {
		return
	}

	clientOptions := options.Client().ApplyURI(mongoDbConnectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	schemaCollection = client.Database("schemaRegistry").Collection("schemas")
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
	serviceSchemasDb, err := schemaCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for serviceSchemasDb.Next(context.TODO()) {
		var elem SchemaRegistry
		err := serviceSchemasDb.Decode(&elem)
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

	if err := serviceSchemasDb.Err(); err != nil {
		log.Fatal(err)
	}

	serviceSchemasDb.Close(context.TODO())

	return nil
}
