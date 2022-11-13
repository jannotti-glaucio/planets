package authenticate

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jannotti-glaucio/planets/core/domains/user/entities"
	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/core/tools/providers/jwt"
	"github.com/jannotti-glaucio/planets/mocks"
)

func TestService(t *testing.T) {
	comm := communication.New()

	useCases := map[string]struct {
		expectedData     string
		expectedResponse communication.Response
		inputData        Dto
		prepare          func(repostitoryMock *mocks.MockIUserRepository, hashMock *mocks.MockIHashProvider, jwtMock *mocks.MockIJwtProvider, loggerMock *mocks.MockILoggerProvider)
	}{
		"success": {
			expectedData: "token",
			expectedResponse: communication.Response{
				Status:  200,
				Code:    comm.Mapping["authenticate_success"].Code,
				Message: comm.Mapping["authenticate_success"].Message,
			},
			inputData: Dto{
				Email:    "user.email@gmail.com",
				Password: "userPassword123",
			},
			prepare: func(repostitoryMock *mocks.MockIUserRepository, hashMock *mocks.MockIHashProvider, jwtMock *mocks.MockIJwtProvider, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(entities.User{UUID: "uuid"}, nil)
				hashMock.EXPECT().Compare(gomock.Any(), gomock.Any()).Return(true)
				jwtMock.EXPECT().CreateToken(gomock.Any()).Return(&jwt.TokenDetails{AccessToken: "token"}, nil)
			},
		},
		"error: on validate": {
			expectedResponse: communication.Response{
				Status:  400,
				Code:    comm.Mapping["validate_failed"].Code,
				Message: comm.Mapping["validate_failed"].Message,
			},
			prepare: func(repostitoryMock *mocks.MockIUserRepository, hashMock *mocks.MockIHashProvider, jwtMock *mocks.MockIJwtProvider, loggerMock *mocks.MockILoggerProvider) {
				loggerMock.EXPECT().Info(gomock.Any(), gomock.Any())
			},
		},
		"error: on repository": {
			expectedData: "",
			expectedResponse: communication.Response{
				Status:  400,
				Code:    comm.Mapping["authenticate_failed"].Code,
				Message: comm.Mapping["authenticate_failed"].Message,
			},
			inputData: Dto{
				Email:    "user.email@gmail.com",
				Password: "userPassword123",
			},
			prepare: func(repostitoryMock *mocks.MockIUserRepository, hashMock *mocks.MockIHashProvider, jwtMock *mocks.MockIJwtProvider, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(entities.User{}, errors.New("error"))
				loggerMock.EXPECT().Error(gomock.Any(), gomock.Any())
			},
		},
		"error: on jwt": {
			expectedData: "",
			expectedResponse: communication.Response{
				Status:  400,
				Code:    comm.Mapping["authenticate_failed"].Code,
				Message: comm.Mapping["authenticate_failed"].Message,
			},
			inputData: Dto{
				Email:    "user.email@gmail.com",
				Password: "userPassword123",
			},
			prepare: func(repostitoryMock *mocks.MockIUserRepository, hashMock *mocks.MockIHashProvider, jwtMock *mocks.MockIJwtProvider, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(entities.User{UUID: "uuid"}, nil)
				hashMock.EXPECT().Compare(gomock.Any(), gomock.Any()).Return(true)
				jwtMock.EXPECT().CreateToken(gomock.Any()).Return(nil, errors.New("error"))
				loggerMock.EXPECT().Error(gomock.Any(), gomock.Any())
			},
		},
		"error: deleted user": {
			expectedData: "",
			expectedResponse: communication.Response{
				Status:  400,
				Code:    comm.Mapping["authenticate_failed"].Code,
				Message: comm.Mapping["authenticate_failed"].Message,
			},
			inputData: Dto{
				Email:    "user.email@gmail.com",
				Password: "userPassword123",
			},
			prepare: func(repostitoryMock *mocks.MockIUserRepository, hashMock *mocks.MockIHashProvider, jwtMock *mocks.MockIJwtProvider, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(entities.User{UUID: ""}, nil)
			},
		},
	}

	for name, useCase := range useCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()
			defer ctrl.Finish()

			repository := mocks.NewMockIUserRepository(ctrl)
			hash := mocks.NewMockIHashProvider(ctrl)
			jwt := mocks.NewMockIJwtProvider(ctrl)
			logger := mocks.NewMockILoggerProvider(ctrl)
			useCase.prepare(repository, hash, jwt, logger)

			service := Service{
				Repository: repository,
				Hash:       hash,
				Jwt:        jwt,
				Logger:     logger,
			}
			data, response := service.Execute(ctx, useCase.inputData)

			if response.Status != useCase.expectedResponse.Status {
				t.Errorf("Expected %d, but got %d", useCase.expectedResponse.Status, response.Status)
			}

			if response.Message != useCase.expectedResponse.Message {
				t.Errorf("Expected %s, but got %s", useCase.expectedResponse.Message, response.Message)
			}

			if data != useCase.expectedData {
				t.Errorf("Expected %s, but got %s", useCase.expectedData, data)
			}
		})
	}
}
