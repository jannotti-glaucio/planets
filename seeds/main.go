package main

import (
	"context"
	"time"

	"github.com/jannotti-glaucio/planets/adapters/persistence/mongodb"
	"github.com/jannotti-glaucio/planets/adapters/persistence/mongodb/repositories/users"
	"github.com/jannotti-glaucio/planets/core/domains/user/services/create"
	"github.com/jannotti-glaucio/planets/core/tools/providers/hash"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`./config.yml`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Second*10))
	defer cancel()

	log := logger.New()

	connection := mongodb.New(ctx)
	if connection.Error != nil {
		log.Error(ctx, "Error connecting to database", connection.Error)
		return
	}

	seedUser(ctx)
}

func seedUser(ctx context.Context) {
	mongoRepo := users.Setup(ctx)
	service := create.Service{Repository: mongoRepo, Hash: hash.New(), Logger: logger.New()}
	dto := create.Dto{
		Email:    "emailteste@gmail.com",
		Password: "123456",
		Name:     "User Test Seed",
	}
	service.Execute(ctx, dto)
}
