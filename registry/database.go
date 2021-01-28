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

// MongoDB client for mongo database
var MongoDB mongoDB

type mongoDBCollection interface {
	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Find(ctx context.Context, filter interface{},
		opts ...*options.FindOptions) (*mongo.Cursor, error)
}

type mongoDB struct {
	collection mongoDBCollection
}

func (m *mongoDB) CreateCollection() {
	mongoDBConnectionString := os.Getenv("MONGODB_CONNECTION_STRING")
	if mongoDBConnectionString == "" {
		log.Fatal("MONGODB_CONNECTION_STRING should not be empty")
	}

	clientOptions := options.Client().ApplyURI(mongoDBConnectionString)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	m.collection = client.Database("schemaRegistry").Collection("schemas")
}

func (m *mongoDB) saveSchema(s *SchemaRegistry) error {
	upsert := true
	updateOptions := options.UpdateOptions{Upsert: &upsert}
	update := bson.M{
		"$set": s,
	}

	filter := bson.D{{Key: "servicename", Value: s.ServiceName}}

	_, err := m.collection.UpdateOne(context.Background(), filter, update, &updateOptions)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoDB) getServiceSchemas(results *[]string) error {
	findOptions := options.Find()
	serviceSchemas, err := m.collection.Find(context.Background(), bson.D{{}}, findOptions)
	if err != nil {
		return err
	}
	for serviceSchemas.Next(context.Background()) {
		var elem SchemaRegistry
		err := serviceSchemas.Decode(&elem)
		if err != nil {
			return err
		}

		value, _ := json.Marshal(elem)
		err = RedisDB.client.Set(ctx, elem.ServiceName, value, 0).Err()
		if err != nil {
			return err
		}

		*results = append(*results, elem.ServiceName)
	}

	if err := serviceSchemas.Err(); err != nil {
		return err
	}

	serviceSchemas.Close(context.Background())

	return nil
}
