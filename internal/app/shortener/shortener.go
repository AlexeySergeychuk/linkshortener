package shortener

import (
	"github.com/AlexeySergeychuk/linkshortener/internal/app/repository"
	"golang.org/x/exp/rand"
)

type ShortenerService interface {
	MakeShortLink(link string) string
	GetFullLink(shortLink string) string
}

type Shortener struct {
	repo repository.Repository
}

func NewShortener(r repository.Repository) *Shortener {
	return &Shortener{repo: r}
} 

// Сохраняет в мапу пару shortLink:longLink
func (s *Shortener) MakeShortLink(link string) string {
	
	// Исключаем отправку разных шортов на один и тот же поданный линк
	if hasCashedLink, shortLink := s.repo.CheckAlreadyHaveShortLink(link); hasCashedLink {
		return makeShortLink(shortLink)
	}

	 shortPath := randomString()
	 s.repo.SaveLinks(shortPath, link)
	
	 return makeShortLink(shortPath)
}

// Возвращаем полный линк
func (s *Shortener) GetFullLink(shortLink string) string {
	return s.repo.GetFullLink(shortLink)
}

// Создает короткую ссылку
func randomString() string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
    rand.Seed(rand.Uint64())

    b := make([]rune, 6)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return "/" + string(b)
}

// Возвращает готовый короткий урл
func makeShortLink(randomstring string) string {
	return "http://localhost:8080" + randomstring
}