package shortener

type Repository interface {
    SaveLinks(shortPath, link string)
	FindByShortLink(shortLink string) string
	FindByFullLink(link string) (bool, string)
}

type ShortLinker interface {
	MakeShortPath(link string) string
}

type Shortener struct {
	repo Repository
	shortLinker ShortLinker
}

func NewShortener(r Repository, s ShortLinker) *Shortener {
	return &Shortener{
		repo: r, 
		shortLinker: s,
	}
} 

// Сохраняет в мапу пару shortLink:longLink
func (s *Shortener) MakeShortLink(link string) string {
	
	// Исключаем отправку разных шортов на один и тот же поданный линк
	if hasCashedLink, shortLink := s.repo.FindByFullLink(link); hasCashedLink {
		return makeShortLink(shortLink)
	}

	 shortPath := s.shortLinker.MakeShortPath(link)
	 s.repo.SaveLinks(shortPath, link)
	
	 return makeShortLink(shortPath)
}

// Возвращаем полный линк
func (s *Shortener) GetFullLink(shortLink string) string {
	return s.repo.FindByShortLink(shortLink)
}


// Возвращает готовый короткий урл
func makeShortLink(randomstring string) string {
	return "http://localhost:8080" + randomstring
}