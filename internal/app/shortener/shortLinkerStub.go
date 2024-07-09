package shortener

import "golang.org/x/exp/rand"

type ShortLinkStub struct {
}

func NewShortLink() *ShortLinkStub {
    return &ShortLinkStub{}
}

// Создает короткую ссылку
func (s *ShortLinkStub) MakeShortPath(link string) string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
    rand.Seed(rand.Uint64())

    b := make([]rune, 6)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return "/" + string(b)
}