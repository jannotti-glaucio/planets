package users

import (
	"context"
	"time"

	"github.com/jannotti-glaucio/planets/adapters/persistence/mongodb"
	"github.com/jannotti-glaucio/planets/adapters/persistence/mongodb/util"
	"github.com/jannotti-glaucio/planets/core/domains/user/entities"
	"github.com/jannotti-glaucio/planets/core/domains/user/repositories"
	"github.com/jannotti-glaucio/planets/core/tools/providers/tracer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//Repository based on user repository from core
type Repository struct {
	Collection *mongo.Collection
}

//Sing repository in interface
var _ repositories.IUserRepository = &Repository{}

//Setup the repository
func Setup(ctx context.Context) *Repository {
	conn := mongodb.New(ctx)
	return &Repository{
		Collection: conn.MongoDB.Collection("users"),
	}
}

//All return all users
func (repo *Repository) All(ctx context.Context) (entities.Users, error) {
	identifierTracer := "mongodb.users.respository.All"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer)
	defer span.Finish()

	var data entities.Users
	cursor, err := repo.Collection.Find(ctx, bson.M{})
	if err != nil {
		return data, err
	}

	for cursor.Next(ctx) {
		row := entities.User{}
		cursor.Decode(&row)
		data = append(data, row)
	}

	return data, nil
}

//Create insert user in DB
func (repo *Repository) Create(ctx context.Context, data entities.User) (entities.User, error) {
	identifierTracer := "mongodb.users.respository.Create"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".filters", Value: data})
	defer span.Finish()

	_, err := repo.Collection.InsertOne(ctx, data)
	if err != nil {
		return entities.User{}, err
	}
	return repo.FindByUUID(ctx, data.UUID)
}

//FindByEmail find user by email
func (repo *Repository) FindByEmail(ctx context.Context, email string) (entities.User, error) {
	identifierTracer := "mongodb.users.respository.FindByEmail"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".email", Value: email})
	defer span.Finish()

	var data entities.User
	err := repo.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//FindByUUID find user by uuid
func (repo *Repository) FindByUUID(ctx context.Context, uuid string) (entities.User, error) {
	identifierTracer := "mongodb.users.respository.FindByUUID"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".uuid", Value: uuid})
	defer span.Finish()

	var data entities.User
	err := repo.Collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//Destroy delete user by uuid
func (repo *Repository) Destroy(ctx context.Context, uuid string) error {
	identifierTracer := "mongodb.users.respository.Destroy"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".uuid", Value: uuid})
	defer span.Finish()

	_, err := repo.Collection.DeleteOne(ctx, bson.M{"_id": uuid})
	if err != nil {
		return err
	}
	return nil
}

//Save update user
func (repo *Repository) Save(ctx context.Context, data entities.User) (entities.User, error) {
	identifierTracer := "mongodb.users.respository.Save"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".data", Value: data})
	defer span.Finish()

	now := time.Now()
	data.UpdatedAt = &now

	primitiveD, err := util.ToUpdate(data)
	if err != nil {
		return entities.User{}, err
	}

	_, err = repo.Collection.UpdateOne(
		ctx,
		bson.M{"_id": data.UUID},
		primitiveD,
	)
	if err != nil {
		return entities.User{}, err
	}
	return repo.FindByUUID(ctx, data.UUID)
}
