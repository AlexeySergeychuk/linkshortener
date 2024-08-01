package repo

import (
	"bufio"
	"encoding/json"
	"os"
)

type Consumer struct {
	file *os.File
	// добавляем reader в Consumer
	reader *bufio.Reader
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file: file,
		// создаём новый Reader
		reader: bufio.NewReader(file),
	}, nil
}

func (c *Consumer) ReadAll() ([]URLdto, error) {
	var urls []URLdto
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}

		var dto URLdto
		if err := json.Unmarshal([]byte(line), &dto); err != nil {
			return nil, err
		}
		urls = append(urls, dto)
	}
	return urls, nil
}
