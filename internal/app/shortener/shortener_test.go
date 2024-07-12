package shortener

import (
	"testing"

	"github.com/AlexeySergeychuk/linkshortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

type MockShortLinker struct {
	mock.Mock
}

func (s *MockRepository) SaveLinks(shortLink, link string) {
	s.Called(shortLink, link)
}

func (s *MockRepository) FindByShortLink(shortLink string) string {
	args := s.Called(shortLink)
	return args.String(0)
}

func (s *MockRepository) FindByFullLink(link string) (bool, string) {
	args := s.Called(link)
	return args.Bool(0), args.String(1)
}

func (s *MockShortLinker) MakeShortPath(link string) string {
	args := s.Called(link)
	return args.String(0)
}

func TestMakeShortLink(t *testing.T) {
	type want struct {
		shortLink string
	}
	tests := []struct {
		name              string
		link              string
		isAlreadyHaveLink bool
		shortPath         string
		want              want
	}{
		{
			name:              "positive case with new link",
			link:              "test.ru",
			isAlreadyHaveLink: false,
			shortPath:         "/rtFgD",
			want: want{
				shortLink: config.FlagShortLinkAddr + "/rtFgD",
			},
		},
		{
			name:              "positive test when link is already in bd",
			link:              "test.ru",
			isAlreadyHaveLink: true,
			shortPath:         "/rtFgD",
			want: want{
				shortLink: config.FlagShortLinkAddr + "/rtFgD",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepository := new(MockRepository)
			mockShortLinker := new(MockShortLinker)

			shortener := NewShortener(mockRepository, mockShortLinker)

			mockRepository.On("FindByFullLink", test.link).Return(test.isAlreadyHaveLink, test.want.shortLink)

			if !test.isAlreadyHaveLink {
				mockRepository.On("SaveLinks", mock.Anything, test.link)
				mockShortLinker.On("MakeShortPath", test.link).Return(test.shortPath)
			}

			shortLink := shortener.MakeShortLink(test.link)
			assert.Equal(t, test.want.shortLink, shortLink)
		})
	}
}
