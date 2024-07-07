package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockShortenerService struct {
	mock.Mock
}

func (s *MockShortenerService) MakeShortLink(link string) string {
	args := s.Called(link)
	return args.String(0)
}

func (s *MockShortenerService) GetFullLink(shortLink string) string {
	args := s.Called(shortLink)
	return args.String(0)
}

func (s *MockShortenerService) CheckAlreadyHaveShortLink(link string) (bool, string) {
	args := s.Called(link)
	return args.Bool(0), args.String(1)
}

func TestCreateLinkHandler(t *testing.T) {
	type want struct {
		code         int
		responseText string
		contentType  string
	}

	tests := []struct {
		name        string
		requestBody string
		want        want
	}{
		{
			name:        "positive test",
			requestBody: "test.ru",
			want: want{
				code:         http.StatusCreated,
				responseText: "http://localhost:8080/rtFgD",
				contentType:  "text/plain",
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
			mockShortener := new(MockShortenerService)

			mockShortener.On("MakeShortLink", test.requestBody).Return(test.want.responseText)

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.requestBody))
			w := httptest.NewRecorder()
			handler := NewHandler(mockShortener)

			// Act
			handler.CreateLinkHandler(w, request)
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
				mockShortener.AssertNotCalled(t, "MakeShortLink", mock.Anything)
			} else {
				mockShortener.AssertExpectations(t)
			}
		})
	}
}

func TestHandler_GetLinkHandler(t *testing.T) {
	type want struct {
		code         int
		headerValue string
	}

	tests := []struct {
		name        string
		path string
		headerName string
		want        want
	}{
		{
			name: "positive test",
			path: "/rtFgD",
			headerName: "Location",
			want: want{
				code: http.StatusTemporaryRedirect,
				headerValue: "test.ru",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Assert
			mockShortener := new(MockShortenerService)
			mockShortener.On("GetFullLink", test.path).Return(test.want.headerValue)

			request := httptest.NewRequest(http.MethodPost, test.path, nil)
			w := httptest.NewRecorder()
			handler := NewHandler(mockShortener)

			// Act
			handler.GetLinkHandler(w, request)
			response := w.Result()
			defer response.Body.Close()

			// Assert
			assert.Equal(t, test.want.code, response.StatusCode)
			assert.Equal(t, test.want.headerValue, w.Header().Get(test.headerName))
		})
	}
}
