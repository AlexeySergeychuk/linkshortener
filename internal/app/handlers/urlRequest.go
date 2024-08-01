package handlers

// easyjson:json
type URLRequest struct {
    URL string `json:"url"`
}

//go:generate easyjson -all urlRequest.go