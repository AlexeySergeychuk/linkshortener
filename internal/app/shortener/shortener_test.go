package shortener

import (
	"testing"

	"github.com/AlexeySergeychuk/linkshortener/internal/app/config"
	"github.com/AlexeySergeychuk/linkshortener/internal/app/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			mockRepository := new(mockRepository)
			mockShortLinker := new(mockShortLinker)
			mockFileProducer := new(mockFileProducer)

			shortener := NewShortener(mockRepository, mockFileProducer, mockShortLinker)

			mockRepository.On("FindByFullLink", test.link).Return(test.isAlreadyHaveLink, test.want.shortLink)

			if !test.isAlreadyHaveLink {
				mockRepository.On("SaveLinks", mock.Anything, test.link)
				mockShortLinker.On("MakeShortPath", test.link).Return(test.shortPath)
				mockFileProducer.On("WriteEvent", repo.URLdto{
					ShortURL: test.shortPath,
					FullURL:  test.link,
				}).Return(nil).Once()
			}

			shortLink := shortener.MakeShortLink(test.link)
			assert.Equal(t, test.want.shortLink, shortLink)

			// Проверка вызова WriteEvent только если ссылка новая
			if !test.isAlreadyHaveLink {
				mockFileProducer.AssertExpectations(t)
			}
		})
	}
}
