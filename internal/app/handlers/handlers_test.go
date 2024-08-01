package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlexeySergeychuk/linkshortener/internal/app/config"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/repo"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/shortener"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRepository struct {
	mock.Mock
}

type mockShortLinker struct {
	mock.Mock
}

type mockFileProducer struct {
	mock.Mock
}

func (s *mockRepository) SaveLinks(shortLink, link string) {
	s.Called(shortLink, link)
}

func (s *mockRepository) FindByShortLink(shortLink string) string {
	args := s.Called(shortLink)
	return args.String(0)
}

func (s *mockRepository) FindByFullLink(link string) (bool, string) {
	args := s.Called(link)
	return args.Bool(0), args.String(1)
}

func (s *mockShortLinker) MakeShortPath(link string) string {
	args := s.Called(link)
	return args.String(0)
}

func (m *mockFileProducer) WriteEvent(event repo.URLdto) error {
	args := m.Called(event)
	return args.Error(0)
}

func TestHandler_CreateShortLinkHandler(t *testing.T) {

	config.FlagShortLinkAddr = "http://localhost:8080"

	type want struct {
		code         int
		responseText string
		contentType  string
	}

	tests := []struct {
		name              string
		requestBody       string
		isAlreadyHaveLink bool
		shortLink         string
		want              want
	}{
		{
			name:              "positive test with new link",
			requestBody:       "test.ru",
			isAlreadyHaveLink: false,
			shortLink:         "/rtFgD",
			want: want{
				code:         http.StatusCreated,
				responseText: config.FlagShortLinkAddr + "/rtFgD",
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name:              "positive test when link is already in bd",
			requestBody:       "test.ru",
			isAlreadyHaveLink: true,
			shortLink:         "/rtFgD1",
			want: want{
				code:         http.StatusCreated,
				responseText: config.FlagShortLinkAddr + "/rtFgD1",
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name:        "no request body",
			requestBody: "",
			want: want{
				code:         http.StatusBadRequest,
				responseText: "",
				contentType:  "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			mockRepository := new(mockRepository)
			mockFileProducer := new(mockFileProducer)
			mockShortLinker := new(mockShortLinker)

			mockRepository.On("FindByFullLink", test.requestBody).Return(test.isAlreadyHaveLink, test.shortLink)

			if !test.isAlreadyHaveLink {
				mockRepository.On("SaveLinks", mock.Anything, test.requestBody)
				mockShortLinker.On("MakeShortPath", test.requestBody).Return(test.shortLink)
				mockFileProducer.On("WriteEvent", mock.Anything).Return(nil)
			}

			shortener := shortener.NewShortener(mockRepository, mockFileProducer, mockShortLinker)
			handler := NewHandler(shortener)

			router := gin.Default()
			router.POST("/", handler.CreateShortLinkHandler)
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.requestBody))
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, request)
			response := w.Result()
			defer response.Body.Close()

			// Assert
			assert.Equal(t, test.want.code, response.StatusCode)
			resBody, err := io.ReadAll(response.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, test.want.responseText, string(resBody))

			// Если requestBody пустой, убедимся что мок не ипользуется
			if test.requestBody == "" {
				mockRepository.AssertNotCalled(t, "FindByFullLink", mock.Anything)
				mockRepository.AssertNotCalled(t, "SaveLinks", mock.Anything, mock.Anything)
				mockShortLinker.AssertNotCalled(t, "MakeShortPath", mock.Anything)
			} else {
				mockRepository.AssertExpectations(t)
				mockShortLinker.AssertExpectations(t)
			}
		})
	}
}

func TestHandler_GetFullLinkHandler(t *testing.T) {
	type want struct {
		code        int
		headerValue string
	}

	tests := []struct {
		name       string
		path       string
		headerName string
		want       want
	}{
		{
			name:       "positive test",
			path:       "/rtFgD",
			headerName: "Location",
			want: want{
				code:        http.StatusTemporaryRedirect,
				headerValue: "/test.ru",
			},
		},
		{
			name:       "negative test when we have no fullLink by shortLink",
			path:       "/rtFgD",
			headerName: "Location",
			want: want{
				code:        http.StatusNotFound,
				headerValue: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Assert
			mockRepository := new(mockRepository)
			mockFileProducer := new(mockFileProducer)
			mockShortLinker := new(mockShortLinker)
			mockRepository.On("FindByShortLink", mock.Anything).Return(test.want.headerValue)

			shortener := shortener.NewShortener(mockRepository, mockFileProducer, mockShortLinker)
			handler := NewHandler(shortener)

			router := gin.Default()
			router.GET("/:id", handler.GetFullLinkHandler)
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, request)
			response := w.Result()
			defer response.Body.Close()

			// Assert
			assert.Equal(t, test.want.code, response.StatusCode)
			assert.Equal(t, test.want.headerValue, w.Header().Get(test.headerName))
		})
	}
}

func TestHandler_GetShortLinkHandler(t *testing.T) {

	config.FlagShortLinkAddr = "http://localhost:8080"

	type want struct {
		code        int
		urlResponse URLResponse
		headerValue string
	}
	tests := []struct {
		name              string
		urlRequest        URLRequest
		isAlreadyHaveLink bool
		shortLink         string
		headerName        string
		want              want
	}{
		{
			name: "positive test",
			urlRequest: URLRequest{
				URL: "https:sometest.ru",
			},
			isAlreadyHaveLink: true,
			shortLink:         "/ggh3t",
			headerName:        "Content-Length",
			want: want{
				code: http.StatusCreated,
				urlResponse: URLResponse{
					Result: config.FlagShortLinkAddr + "/ggh3t",
				},
				headerValue: "40",
			},
		},
		{
			name: "url is empty",
			urlRequest: URLRequest{
				URL: "",
			},
			isAlreadyHaveLink: true,
			shortLink:         "",
			headerName:        "Content-Length",
			want: want{
				code: http.StatusBadRequest,
				urlResponse: URLResponse{
					Result: "",
				},
				headerValue: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Assert
			mockRepository := new(mockRepository)
			mockFileProducer := new(mockFileProducer)
			mockShortLinker := new(mockShortLinker)
			mockRepository.On("FindByFullLink", mock.Anything).Return(test.isAlreadyHaveLink, test.shortLink)

			shortener := shortener.NewShortener(mockRepository, mockFileProducer, mockShortLinker)
			handler := NewHandler(shortener)

			router := gin.Default()
			router.POST("/api/shorten", handler.GetShortLinkHandler)

			body, err := easyjson.Marshal(test.urlRequest)
			require.NoError(t, err, "failed to read response body")
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, request)
			response := w.Result()
			defer response.Body.Close()

			// Assert
			assert.Equal(t, test.want.code, response.StatusCode)
			if response.StatusCode == http.StatusOK {
				assert.Equal(t, test.want.headerValue, response.Header.Get(test.headerName))

				var resBody URLResponse
				err = easyjson.UnmarshalFromReader(response.Body, &resBody)
				require.NoError(t, err)
				assert.Equal(t, test.want.urlResponse, resBody)
			}
		})
	}
}
