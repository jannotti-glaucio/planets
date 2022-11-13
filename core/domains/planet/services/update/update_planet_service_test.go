package update

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jannotti-glaucio/planets/core/domains/planet/entities"
	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/mocks"
)

func TestService(t *testing.T) {
	comm := communication.New()

	expectedData := entities.Planet{
		UUID:    "id",
		Name:    "name",
		Terrain: "terrain",
		Climate: "climate",
	}

	useCases := map[string]struct {
		expectedResponse communication.Response
		expectedData     entities.Planet
		inputData        Dto
		prepare          func(repostitoryMock *mocks.MockIPlanetRepository, loggerMock *mocks.MockILoggerProvider, cacheMock *mocks.MockICacheProvider)
	}{
		"success": {
			inputData: Dto{
				UUID:    "uuid",
				Name:    "name",
				Terrain: "terrain",
				Climate: "climate",
			},
			expectedData: expectedData,
			expectedResponse: communication.Response{
				Status:  200,
				Code:    comm.Mapping["success"].Code,
				Message: comm.Mapping["success"].Message,
			},
			prepare: func(repostitoryMock *mocks.MockIPlanetRepository, loggerMock *mocks.MockILoggerProvider, cacheMock *mocks.MockICacheProvider) {
				repostitoryMock.EXPECT().FindByUUID(gomock.Any(), gomock.Any()).Return(expectedData, nil)
				repostitoryMock.EXPECT().FindByName(gomock.Any(), gomock.Any()).Return(entities.Planet{}, nil)
				cacheMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				repostitoryMock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(expectedData, nil)
			},
		},
		"error: on repository FindByUUID": {
			expectedResponse: communication.Response{
				Status:  500,
				Code:    comm.Mapping["error_update"].Code,
				Message: comm.Mapping["error_update"].Message,
			},
			inputData: Dto{
				UUID:    "uuid",
				Name:    "name",
				Terrain: "terrain",
				Climate: "climate",
			},
			prepare: func(repostitoryMock *mocks.MockIPlanetRepository, loggerMock *mocks.MockILoggerProvider, cacheMock *mocks.MockICacheProvider) {
				repostitoryMock.EXPECT().FindByUUID(gomock.Any(), gomock.Any()).Return(entities.Planet{}, errors.New("error"))
				loggerMock.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
		},
		"error: on validation": {
			expectedResponse: communication.Response{
				Status:  400,
				Code:    comm.Mapping["validate_failed"].Code,
				Message: comm.Mapping["validate_failed"].Message,
			},
			inputData: Dto{
				UUID:    "uuid",
				Name:    "name",
				Terrain: "terrain",
				Climate: "climate",
			},
			prepare: func(repostitoryMock *mocks.MockIPlanetRepository, loggerMock *mocks.MockILoggerProvider, cacheMock *mocks.MockICacheProvider) {
				repostitoryMock.EXPECT().FindByUUID(gomock.Any(), gomock.Any()).Return(entities.Planet{}, nil)
				repostitoryMock.EXPECT().FindByName(gomock.Any(), gomock.Any()).Return(entities.Planet{UUID: "uuid"}, errors.New("error"))
				loggerMock.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any())
				loggerMock.EXPECT().Info(gomock.Any(), gomock.Any())
			},
		},
		"error: on repository Save": {
			expectedResponse: communication.Response{
				Status:  500,
				Code:    comm.Mapping["error_update"].Code,
				Message: comm.Mapping["error_update"].Message,
			},
			inputData: Dto{
				UUID:    "uuid",
				Name:    "name",
				Terrain: "terrain",
				Climate: "climate",
			},
			prepare: func(repostitoryMock *mocks.MockIPlanetRepository, loggerMock *mocks.MockILoggerProvider, cacheMock *mocks.MockICacheProvider) {
				repostitoryMock.EXPECT().FindByUUID(gomock.Any(), gomock.Any()).Return(expectedData, nil)
				repostitoryMock.EXPECT().FindByName(gomock.Any(), gomock.Any()).Return(entities.Planet{}, nil)
				cacheMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				repostitoryMock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(entities.Planet{}, errors.New("error"))
				loggerMock.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
		},
	}

	for name, useCase := range useCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()
			defer ctrl.Finish()

			repository := mocks.NewMockIPlanetRepository(ctrl)
			logger := mocks.NewMockILoggerProvider(ctrl)
			client := mocks.NewMockIHTTPClientProvider(ctrl)
			cache := mocks.NewMockICacheProvider(ctrl)
			useCase.prepare(repository, logger, cache)

			service := Service{
				Repository: repository,
				Logger:     logger,
				HTTPClient: client,
				Cache:      cache,
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
