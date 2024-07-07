package shortener

import (
	"github.com/AlexeySergeychuk/linkshortener/internal/app/repository"
	"golang.org/x/exp/rand"
)

// Сохраняет в мапу пару shortLink:longLink
func MakeShortLink(link string) string {
	
	// Исключаем отправку разных шортов на один и тот же поданный линк
	if hasCashedLink, shortLink := repository.CheckAlreadyHaveShortLink(link); hasCashedLink {
		return makeShortLink(shortLink)
	}

	 shortPath := randomString()
	 repository.SaveLinks(shortPath, link)
	
	 return makeShortLink(shortPath)
}

// Возвращаем полный линк
func GetFullLink(shortLink string) string {
	return repository.GetFullLink(shortLink)
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