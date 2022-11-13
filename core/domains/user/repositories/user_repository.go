package repositories

import (
	"context"

	"github.com/jannotti-glaucio/planets/core/domains/user/entities"
)

//IUserRepository ...
type IUserRepository interface {
	All(ctx context.Context) (entities.Users, error)
	Create(ctx context.Context, data entities.User) (entities.User, error)
	FindByEmail(ctx context.Context, email string) (entities.User, error)
	FindByUUID(ctx context.Context, uuid string) (entities.User, error)
	Destroy(ctx context.Context, uuid string) error
	Save(ctx context.Context, data entities.User) (entities.User, error)
}
