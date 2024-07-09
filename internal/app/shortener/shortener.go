package shortener

import (
	"github.com/AlexeySergeychuk/linkshortener/internal/app/repository"
)

type ShortenerService interface {
	MakeShortLink(link string) string
	GetFullLink(shortLink string) string
}

type Shortener struct {
	repo repository.Repository
	shortLinker ShortLinker
}

func NewShortener(r repository.Repository, s ShortLinker) *Shortener {
	return &Shortener{
		repo: r, 
		shortLinker: s,
	}
} 

// Сохраняет в мапу пару shortLink:longLink
func (s *Shortener) MakeShortLink(link string) string {
	
	// Исключаем отправку разных шортов на один и тот же поданный линк
	if hasCashedLink, shortLink := s.repo.CheckAlreadyHaveShortLink(link); hasCashedLink {
		return makeShortLink(shortLink)
	}

	 shortPath := s.shortLinker.RandomString(link)
	 s.repo.SaveLinks(shortPath, link)
	
	 return makeShortLink(shortPath)
}

// Возвращаем полный линк
func (s *Shortener) GetFullLink(shortLink string) string {
	return s.repo.GetFullLink(shortLink)
}


// Возвращает готовый короткий урл
func makeShortLink(randomstring string) string {
	return "http://localhost:8080" + randomstring
}