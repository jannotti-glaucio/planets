package create

import (
	"context"
	"encoding/json"

	"github.com/jannotti-glaucio/planets/core/domains/user/entities"
	"github.com/jannotti-glaucio/planets/core/domains/user/repositories"
	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/core/tools/providers/hash"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
	"github.com/jannotti-glaucio/planets/core/tools/providers/tracer"
	"github.com/jannotti-glaucio/planets/core/tools/validations"
)

//Dto object receiver
type Dto struct {
	Name     string `json:"name" validate:"isRequired"`
	Email    string `json:"email" validate:"isRequired|isEmail"`
	Password string `json:"password" validate:"isPassword"`
}

//Service ...
type Service struct {
	Repository repositories.IUserRepository
	Hash       hash.IHashProvider
	Logger     logger.ILoggerProvider
}

//Execute Serviço responsável pela inserção de registros
func (service *Service) Execute(ctx context.Context, dto Dto) (created entities.User, response communication.Response) {
	identifierTracer := "create.user.service"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".dto", Value: dto})
	defer span.Finish()

	response.Fields = validations.ValidateStruct(&dto, "")
	comm := communication.New()

	userByEmail, err := service.Repository.FindByEmail(ctx, dto.Email)
	if err != nil {
		service.Logger.Info(ctx, "domain.user.service.create.create_user_service.Repository.FindByEmail", err)
	}

	//Check e-mail in use
	if userByEmail.UUID != "" {
		response.Fields = append(response.Fields, comm.Fields("email", "already_exists"))
	}

	if len(response.Fields) > 0 {
		service.Logger.Info(ctx, "domain.user.service.create.create_user_service.ValidationError")
		resp := comm.Response(400, "validate_failed")
		resp.Fields = response.Fields
		response = resp
		return
	}

	//Apply security hash in password
	dto.Password = service.Hash.Create(dto.Password)

	User := entities.User{}
	user := User.New()

	//Mergin entity and dto
	toMerge, _ := json.Marshal(dto)
	json.Unmarshal(toMerge, &user)

	created, err = service.Repository.Create(ctx, *user)

	if err != nil {
		service.Logger.Error(ctx, "domain.user.service.create.create_user_service.Repository.Create", err)
		response = comm.Response(500, "error_create")
		return
	}

	response = comm.Response(200, "success")
	return
}
