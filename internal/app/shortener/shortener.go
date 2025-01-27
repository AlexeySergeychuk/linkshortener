package shortener

import "github.com/AlexeySergeychuk/linkshortener/internal/app/config"

type Repository interface {
	SaveLinks(shortPath, link string)
	FindByShortLink(shortLink string) string
	FindByFullLink(link string) (bool, string)
}

type ShortLinker interface {
	MakeShortPath(link string) string
}

type Shorten struct {
	repo        Repository
	shortLinker ShortLinker
}

func NewShortener(r Repository, s ShortLinker) *Shorten {
	return &Shorten{
		repo:        r,
		shortLinker: s,
	}
}

// Сохраняет в мапу пару shortLink:longLink
func (s *Shorten) MakeShortLink(link string) string {

	// Исключаем отправку разных шортов на один и тот же поданный линк
	if hasSavedLink, shortLink := s.repo.FindByFullLink(link); hasSavedLink {
		return makeShortLink(shortLink)
	}

	shortPath := s.shortLinker.MakeShortPath(link)
	s.repo.SaveLinks(shortPath, link)

	return makeShortLink(shortPath)
}

// Возвращаем полный линк
func (s *Shorten) GetFullLink(shortLink string) string {
	return s.repo.FindByShortLink(shortLink)
}

// Возвращает готовый короткий урл
func makeShortLink(randomstring string) string {
	return config.FlagShortLinkAddr + randomstring
}
