package films

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jannotti-glaucio/planets/mocks"
)

func TestService(t *testing.T) {
	useCases := map[string]struct {
		expectedData int
		inputData    string
		prepare      func(loggerMock *mocks.MockILoggerProvider, clientMock *mocks.MockIHTTPClientProvider, cacheMock *mocks.MockICacheProvider)
	}{
		"success": {
			inputData:    "planetName",
			expectedData: 2,
			prepare: func(loggerMock *mocks.MockILoggerProvider, clientMock *mocks.MockIHTTPClientProvider, cacheMock *mocks.MockICacheProvider) {
				cacheMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
				bodyResponse := ioutil.NopCloser(bytes.NewReader([]byte(`{"results": [{"films": ["film1", "film2"]}]}`)))
				clientMock.EXPECT().Get(gomock.Any()).Times(1).Return(&http.Response{Body: bodyResponse}, nil)
				cacheMock.EXPECT().WithExpiration(gomock.Any()).Times(1).Return(cacheMock)
				cacheMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
			},
		},
		"success: from cache": {
			inputData:    "planetName",
			expectedData: 6,
			prepare: func(loggerMock *mocks.MockILoggerProvider, clientMock *mocks.MockIHTTPClientProvider, cacheMock *mocks.MockICacheProvider) {
				cacheMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).SetArg(2, 6).Return(nil)
			},
		},
		"error: on Client Request": {
			inputData: "planetName",
			prepare: func(loggerMock *mocks.MockILoggerProvider, clientMock *mocks.MockIHTTPClientProvider, cacheMock *mocks.MockICacheProvider) {
				cacheMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
				clientMock.EXPECT().Get(gomock.Any()).Times(1).Return(nil, errors.New("error"))
				loggerMock.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
			},
		},
		"error: on Response Body": {
			inputData: "planetName",
			prepare: func(loggerMock *mocks.MockILoggerProvider, clientMock *mocks.MockIHTTPClientProvider, cacheMock *mocks.MockICacheProvider) {
				cacheMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
				bodyResponse := ioutil.NopCloser(bytes.NewReader([]byte(`{invalidJson}`)))
				clientMock.EXPECT().Get(gomock.Any()).Times(1).Return(&http.Response{Body: bodyResponse}, nil)
				loggerMock.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
			},
		},
		"success: zero results": {
			inputData: "planetName",
			prepare: func(loggerMock *mocks.MockILoggerProvider, clientMock *mocks.MockIHTTPClientProvider, cacheMock *mocks.MockICacheProvider) {
				cacheMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
				bodyResponse := ioutil.NopCloser(bytes.NewReader([]byte(`{"results": []}`)))
				clientMock.EXPECT().Get(gomock.Any()).Times(1).Return(&http.Response{Body: bodyResponse}, nil)
			},
		},
	}

	for name, useCase := range useCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()
			defer ctrl.Finish()

			logger := mocks.NewMockILoggerProvider(ctrl)
			client := mocks.NewMockIHTTPClientProvider(ctrl)
			cache := mocks.NewMockICacheProvider(ctrl)
			useCase.prepare(logger, client, cache)

			service := Service{
				Logger:     logger,
				HTTPClient: client,
				Cache:      cache,
			}
			data := service.Execute(ctx, useCase.inputData)

			if data != useCase.expectedData {
				t.Errorf("Expected %d, but got %d", useCase.expectedData, data)
			}
		})
	}
}
