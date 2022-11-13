package planets

import (
	"context"
	"time"

	"github.com/jannotti-glaucio/planets/adapters/persistence/mongodb"
	"github.com/jannotti-glaucio/planets/adapters/persistence/mongodb/util"
	"github.com/jannotti-glaucio/planets/core/domains/planet/entities"
	"github.com/jannotti-glaucio/planets/core/domains/planet/repositories"
	"github.com/jannotti-glaucio/planets/core/tools/providers/tracer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//Repository based on planet repository from core
type Repository struct {
	Collection *mongo.Collection
}

//Sing repository in interface
var _ repositories.IPlanetRepository = &Repository{}

//Setup ...
func Setup(ctx context.Context) *Repository {
	conn := mongodb.New(ctx)
	return &Repository{
		Collection: conn.MongoDB.Collection("planets"),
	}
}

//All return all planets
func (repo *Repository) All(ctx context.Context, filters *map[string]interface{}) (entities.Planets, error) {
	identifierTracer := "mongodb.planets.respository.All"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".filters", Value: filters})
	defer span.Finish()

	var data entities.Planets

	cursor, err := repo.Collection.Find(ctx, filters)
	if err != nil {
		return data, err
	}

	for cursor.Next(ctx) {
		row := entities.Planet{}
		cursor.Decode(&row)
		data = append(data, row)
	}

	return data, nil
}

//Create insert planet in DB
func (repo *Repository) Create(ctx context.Context, data entities.Planet) (entities.Planet, error) {
	identifierTracer := "mongodb.planets.respository.Create"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".data", Value: data})
	defer span.Finish()

	_, err := repo.Collection.InsertOne(ctx, data)
	if err != nil {
		return entities.Planet{}, err
	}
	return repo.FindByUUID(ctx, data.UUID)
}

//FindByUUID find planet by uuid
func (repo *Repository) FindByUUID(ctx context.Context, uuid string) (entities.Planet, error) {
	identifierTracer := "mongodb.planets.respository.FindByUUID"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".uuid", Value: uuid})
	defer span.Finish()

	var data entities.Planet
	err := repo.Collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//FindByName find planet by uuid
func (repo *Repository) FindByName(ctx context.Context, name string) (entities.Planet, error) {
	identifierTracer := "mongodb.planets.respository.FindByName"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".name", Value: name})
	defer span.Finish()

	var data entities.Planet
	err := repo.Collection.FindOne(ctx, bson.M{"name": name}).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//Destroy delete planet by uuid
func (repo *Repository) Destroy(ctx context.Context, uuid string) error {
	identifierTracer := "mongodb.planets.respository.Destroy"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".uuid", Value: uuid})
	defer span.Finish()

	_, err := repo.Collection.DeleteOne(ctx, bson.M{"_id": uuid})
	if err != nil {
		return err
	}
	return nil
}

//Save update planet
func (repo *Repository) Save(ctx context.Context, data entities.Planet) (entities.Planet, error) {
	identifierTracer := "mongodb.planets.respository.Save"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".data", Value: data})
	defer span.Finish()

	now := time.Now()
	data.UpdatedAt = &now

	primitiveD, err := util.ToUpdate(data)
	if err != nil {
		return entities.Planet{}, err
	}

	_, err = repo.Collection.UpdateOne(
		ctx,
		bson.M{"_id": data.UUID},
		primitiveD,
	)
	if err != nil {
		return entities.Planet{}, err
	}
	return repo.FindByUUID(ctx, data.UUID)
}
