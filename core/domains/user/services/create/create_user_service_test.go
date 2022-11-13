package create

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jannotti-glaucio/planets/core/domains/user/entities"
	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/mocks"
)

func TestService(t *testing.T) {
	comm := communication.New()

	expectedData := entities.User{
		UUID:  "userID",
		Name:  "userName",
		Email: "userEmail",
	}

	useCases := map[string]struct {
		expectedData     entities.User
		expectedResponse communication.Response
		inputData        Dto
		prepare          func(repostitoryMock *mocks.MockIUserRepository, hashMock *mocks.MockIHashProvider, loggerMock *mocks.MockILoggerProvider)
	}{
		"success": {
			expectedData: expectedData,
			expectedResponse: communication.Response{
				Status:  200,
				Code:    comm.Mapping["success"].Code,
				Message: comm.Mapping["success"].Message,
			},
			inputData: Dto{
				Name:     "userName",
				Email:    "user.email@gmail.com",
				Password: "userPassword123",
			},
			prepare: func(repostitoryMock *mocks.MockIUserRepository, hashMock *mocks.MockIHashProvider, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(entities.User{}, nil)
				hashMock.EXPECT().Create(gomock.Any()).Return("hash")
				repostitoryMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(expectedData, nil)
			},
		},
		"error: on validation": {
			expectedResponse: communication.Response{
				Status:  400,
				Code:    comm.Mapping["validate_failed"].Code,
				Message: comm.Mapping["validate_failed"].Message,
			},
			inputData: Dto{
				Name:     "userName",
				Email:    "user.email@gmail.com",
				Password: "userPassword123",
			},
			prepare: func(repostitoryMock *mocks.MockIUserRepository, hashMock *mocks.MockIHashProvider, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(entities.User{UUID: "OuterUserID"}, errors.New("error"))
				loggerMock.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any())
				loggerMock.EXPECT().Info(gomock.Any(), gomock.Any())
			},
		},
		"error: on repository Create": {
			expectedResponse: communication.Response{
				Status:  500,
				Code:    comm.Mapping["error_create"].Code,
				Message: comm.Mapping["error_create"].Message,
			},
			inputData: Dto{
				Name:     "userName",
				Email:    "user.email@gmail.com",
				Password: "userPassword123",
			},
			prepare: func(repostitoryMock *mocks.MockIUserRepository, hashMock *mocks.MockIHashProvider, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(entities.User{}, nil)
				hashMock.EXPECT().Create(gomock.Any()).Return("hash")
				repostitoryMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(entities.User{}, errors.New("error"))
				loggerMock.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
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
			logger := mocks.NewMockILoggerProvider(ctrl)
			useCase.prepare(repository, hash, logger)

			service := Service{
				Repository: repository,
				Hash:       hash,
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
				t.Errorf("Expected %q, but got %q", useCase.expectedData, data)
			}
		})
	}
}
