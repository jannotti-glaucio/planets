package mongodb

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var singletonConnection *Connection = nil

//New initialize a new connection
func New(ctx context.Context) *Connection {
	if singletonConnection == nil {
		mongoDB, err := startMongoConnection(ctx)
		singletonConnection = &Connection{
			MongoDB: mongoDB,
			Error:   err,
		}
	}
	return singletonConnection
}

//Connection database connection
type Connection struct {
	Error   error
	MongoDB *mongo.Database
}

//startMongoConnection start a new connection
func startMongoConnection(ctx context.Context) (*mongo.Database, error) {

	connectionURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		viper.GetString("mongo.user"),
		viper.GetString("mongo.password"),
		viper.GetString("mongo.host"),
		viper.GetString("mongo.port"),
		viper.GetString("mongo.database"),
	)

	connection := options.Client().ApplyURI(connectionURI)
	client, err := mongo.Connect(ctx, connection)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(viper.GetString("mongo.database")), nil
}
