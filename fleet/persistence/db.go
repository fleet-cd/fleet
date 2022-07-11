package persistence

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/tgs266/fleet/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func Connect(config config.Config) {
	log.Info().Msg("connecting to mongodb")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoDB.URI))
	if err != nil {
		panic(err)
	}

	db = client.Database(config.MongoDB.Database)
	createCollection("ships")
	createCollection("products")
	createCollection("cargo")
}

func createIndex(name string, opts *options.IndexOptions) mongo.IndexModel {
	return mongo.IndexModel{
		Keys: bson.M{
			name: 1,
		}, Options: &options.IndexOptions{},
	}
}

func createCollection(collection string) {
	createdAtIndex := createIndex("createdAt", options.Index())
	modifiedAtIndex := createIndex("modifiedAt", options.Index())
	t := new(bool)
	*t = true
	frnIndex := createIndex("frn", &options.IndexOptions{Unique: t})

	col := db.Collection(collection)
	col.Indexes().CreateOne(context.TODO(), createdAtIndex, options.CreateIndexes())
	col.Indexes().CreateOne(context.TODO(), modifiedAtIndex, options.CreateIndexes())
	col.Indexes().CreateOne(context.TODO(), frnIndex, options.CreateIndexes())
}

func GetCollection(collection string) *mongo.Collection {
	col := db.Collection(collection)
	col.Indexes().CreateOne(context.TODO(), mongo.IndexModel{})
	return col
}

func FindOneByFrn[T any](ctx context.Context, collection string, frn string) (T, error) {
	col := GetCollection(collection)
	var val T
	err := col.FindOne(ctx, bson.M{"frn": frn}, &options.FindOneOptions{}).Decode(&val)
	return val, err
}

func InsertOneToCollection(ctx context.Context, collection string, object any) error {
	col := GetCollection(collection)
	_, err := col.InsertOne(ctx, object)
	return err
}

func DecodeCursor[T any](ctx context.Context, cur *mongo.Cursor) ([]T, error) {
	results := []T{}
	for cur.Next(ctx) {
		var elem T
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	return results, nil
}
