package repo

type URLdto struct {
	UUID     uint   `json:"uuid"`
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}