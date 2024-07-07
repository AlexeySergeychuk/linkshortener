package repository

// var mapStorage = map[string]string{}

type Repository interface {
    SaveLinks(shortPath, link string)
	GetFullLink(shortLink string) string
	CheckAlreadyHaveShortLink(link string) (bool, string)
}

type Repo struct {
	bdStub map[string]string
}

func NewRepo(bd map[string]string) *Repo {
	return &Repo{bdStub: bd}
}

// Сохранить в мапу оба линка
func (r *Repo) SaveLinks(shortLink, link string) {
	r.bdStub[shortLink] = link
}

// Вернуть full link
func (r *Repo) GetFullLink(shortLink string) string {
	return r.bdStub[shortLink]
}

// Проверить, не был ли ранее записан этот шорт линк
func (r *Repo) CheckAlreadyHaveShortLink(link string) (bool, string) {
	for k, v := range r.bdStub {
		if v == link {
			return true, k
		}
	}		
	return false, ""	
}