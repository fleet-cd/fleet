package persistence

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/tgs266/fleet/config"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
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
}

func getCollection(collection string) *mongo.Collection {
	col := db.Collection(collection)
	return col
}

func CreateShip(ctx context.Context, ship entities.ShipEntity) error {
	col := getCollection("ships")
	_, err := col.InsertOne(ctx, ship)
	return err
}

func GetShip(ctx context.Context, frn string) (entities.ShipEntity, error) {
	col := getCollection("ships")
	var ship entities.ShipEntity
	err := col.FindOne(ctx, bson.M{"frn": frn}, &options.FindOneOptions{}).Decode(&ship)
	return ship, err
}

func ListShips(ctx context.Context) ([]entities.ShipEntity, error) {
	col := getCollection("ships")
	cur, err := col.Find(ctx, bson.M{}, &options.FindOptions{})
	if err != nil {
		return nil, err
	}
	results := []entities.ShipEntity{}
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem entities.ShipEntity
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)

	}
	return results, nil
}

// func GetShip(ctx context.Context, frn frn.String) (entities.ShipEntity, error) {
// 	col := getCollection("ships")
// 	var ship entities.ShipEntity
// 	if err := col.FindOne(ctx, bson.M{"frn": frn}).Decode(&ship); err != nil {
// 		return entities.ShipEntity{}, err
// 	}
// 	return ship, nil
// }
