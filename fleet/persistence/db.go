package persistence

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tgs266/fleet/config"
	"github.com/tgs266/fleet/fleet/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var cfg config.Config

func Connect(config *config.Config) error {
	if db != nil {
		return nil
	}
	log.Info().Msg("connecting to mongodb")
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	cfg = utils.OrDefault(config, cfg)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoDB.URI).SetConnectTimeout(1*time.Second))
	if err != nil {
		log.Logger.Err(err).Msg("could not connect to database")
		return err
	}

	db = client.Database(config.MongoDB.Database)
	createCollection("permissions", true, "name")
	createCollection("groups", false, "name")
	createCollection("users", true, "email")
	createCollection("ships", true)
	createCollection("products", true)
	createCollection("cargo", true)
	return nil
}

func createIndex(name string, opts *options.IndexOptions) mongo.IndexModel {
	return mongo.IndexModel{
		Keys: bson.M{
			name: 1,
		}, Options: opts,
	}
}

func createCollection(collection string, useFrn bool, uniques ...string) {
	createdAtIndex := createIndex("createdAt", options.Index())
	modifiedAtIndex := createIndex("modifiedAt", options.Index())
	t := new(bool)
	*t = true
	frnIndex := createIndex("frn", &options.IndexOptions{Unique: t})

	col := db.Collection(collection)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	col.Indexes().CreateOne(ctx, createdAtIndex, options.CreateIndexes())

	ctx, cancel = context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	col.Indexes().CreateOne(ctx, modifiedAtIndex, options.CreateIndexes())

	if useFrn {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Second*2)
		defer cancel()
		col.Indexes().CreateOne(ctx, frnIndex, options.CreateIndexes())
	}

	for _, s := range uniques {
		index := createIndex(s, &options.IndexOptions{Unique: t})
		ctx, cancel = context.WithTimeout(context.TODO(), time.Second*2)
		defer cancel()
		col.Indexes().CreateOne(ctx, index, options.CreateIndexes())
	}
}

func Ping() bool {
	if db == nil {
		err := Connect(nil)
		if err != nil {
			return false
		}
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	err := db.Client().Ping(ctx, nil)
	if err != nil {
		return false
	}
	return true
}

func GetCollection(collection string) (*mongo.Collection, error) {
	err := Connect(nil)
	if err != nil {
		return nil, err
	}
	col := db.Collection(collection)
	col.Indexes().CreateOne(context.TODO(), mongo.IndexModel{})
	return col, nil
}

func FindOneByFrn[T any](ctx context.Context, collection string, frn string) (T, error) {
	var val T
	col, err := GetCollection(collection)
	if err != nil {
		return val, err
	}
	err = col.FindOne(ctx, bson.M{"frn": frn}, &options.FindOneOptions{}).Decode(&val)
	return val, err
}

func DeleteOneByFrn[T any](ctx context.Context, collection string, frn string) error {
	col, err := GetCollection(collection)
	if err != nil {
		return err
	}
	_, err = col.DeleteOne(ctx, bson.M{"frn": frn}, &options.DeleteOptions{})
	return err
}

func FindOne[T any](ctx context.Context, collection string, filter bson.M) (T, error) {
	var val T
	col, err := GetCollection(collection)
	if err != nil {
		return val, err
	}
	err = col.FindOne(ctx, filter, &options.FindOneOptions{}).Decode(&val)
	return val, err
}

func Count(ctx context.Context, collection string) (int64, error) {
	col, err := GetCollection(collection)
	if err != nil {
		return -1, err
	}
	return col.CountDocuments(ctx, bson.M{}, &options.CountOptions{})
}

func List[T any](ctx context.Context, collection string, opts *options.FindOptions) ([]T, error) {
	col, err := GetCollection(collection)
	if err != nil {
		return nil, err
	}
	cur, err := col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	return DecodeCursor[T](ctx, cur)
}

func InsertOneToCollection(ctx context.Context, collection string, object any) error {
	col, err := GetCollection(collection)
	if err != nil {
		return err
	}
	_, err = col.InsertOne(ctx, object)
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
