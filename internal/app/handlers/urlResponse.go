package handlers

// easyjson:json
type URLResponse struct {
	Result string `json:"result"`
}

//go:generate easyjson -all urlResponse.go
