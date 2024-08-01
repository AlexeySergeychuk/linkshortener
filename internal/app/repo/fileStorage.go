package repo

type URLdto struct {
	UUID     uint   `json:"uuid"`
	ShortUrl string `json:"short_url"`
	FullUrl  string `json:"original_url"`
}