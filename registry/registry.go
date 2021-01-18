package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SchemaRegistry to add a schema to the registry
type SchemaRegistry struct {
	ServiceName string
	ServiceURL  string
	TypeDefs    string
}

var redisDB = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var ctx = context.Background()

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

// ValidateSchema validates the graphql schema
func (s *SchemaRegistry) ValidateSchema() error {
	if s.TypeDefs == "" {
		return fmt.Errorf("typedefs should not be empty")
	}
	_, err := gqlparser.LoadSchema(&ast.Source{Name: s.ServiceName, Input: s.TypeDefs, BuiltIn: false})
	if err != nil {
		return err
	}
	return nil
}

// Save the schema in redis
func (s *SchemaRegistry) Save() error {
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

	value, _ := json.Marshal(s)
	err = redisDB.Set(ctx, s.ServiceName, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetServiceSchema get service schema from redis
func GetServiceSchema(service string) (*SchemaRegistry, error) {
	val2, err := redisDB.Get(ctx, service).Result()
	if err == redis.Nil {
		return &SchemaRegistry{}, err
	} else if err != nil {
		return &SchemaRegistry{}, err
	} else {
		var serviceSchema SchemaRegistry
		err := json.Unmarshal([]byte(val2), &serviceSchema)
		if err != nil {
			return &SchemaRegistry{}, err
		}
		return &serviceSchema, nil
	}
}

// GetAllServices returns all services names
func GetAllServices() (*[]string, error) {
	var cursor uint64
	serviceSchemasRd, _, err := redisDB.Scan(ctx, cursor, "*", 100).Result()
	if err != nil {
		return &[]string{}, err
	}

	if len(serviceSchemasRd) == 0 {
		var results []string

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
				return nil, err
			}

			results = append(results, elem.ServiceName)
		}

		if err := serviceSchemasDb.Err(); err != nil {
			log.Fatal(err)
		}

		serviceSchemasDb.Close(context.TODO())

		return &results, nil
	}
	return &serviceSchemasRd, nil
}
