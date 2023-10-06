package files

import (
	"encoding/json"
	Event "github.com/Buff2out/shurle/internal/app/api/shortener"
	"os"
)

type Producer struct {
	file *os.File // Файл для записи
}

type Consumer struct {
	file *os.File // Файл для чтения
}

func NewProducer(filename string) (*Producer, error) {
	// открываем файл для записи в конец
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{file: file}, nil
}

func (p *Producer) Close() error {
	// закрываем файл
	return p.file.Close()
}

func NewConsumer(filename string) (*Consumer, error) {
	// открываем файл для чтения
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{file: file}, nil
}

func (c *Consumer) Close() error {
	// закрываем файл
	return c.file.Close()
}

func (p *Producer) WriteEvent(event *Event.ShURLFile) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}
	// добавляем перенос строки
	data = append(data, '\n')

	_, err = p.file.Write(data)
	return err
}
