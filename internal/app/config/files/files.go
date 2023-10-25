package files

import (
	"bufio"
	"encoding/json"
	"fmt"
	Event "github.com/Buff2out/shurle/internal/app/api/shortener"
	"go.uber.org/zap"
	"os"
)

type Producer struct {
	file *os.File // Файл для записи
	// добавляем Writer в Producer
	writer *bufio.Writer
}

type Consumer struct {
	file *os.File // Файл для чтения
	// заменяем Reader на Scanner
	scanner *bufio.Scanner
}

func NewProducer(filename string, sugar *zap.SugaredLogger) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		sugar.Infow("Is file not created?? WTF ", "errorMsg ", err)
		return nil, err
	}

	return &Producer{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func NewConsumer(filename string, sugar *zap.SugaredLogger) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		sugar.Infow("Is file not created?? WTF ", "errorMsg ", err)
		return nil, err
	}

	return &Consumer{
		file: file,
		// создаём новый scanner
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*Event.ShURLFile, error) {
	// одиночное сканирование до следующей строки
	if !c.scanner.Scan() {
		return nil, fmt.Errorf("КОНЕЦ СТРОКИ")
	}
	// читаем данные из scanner
	data := c.scanner.Bytes()

	event := Event.ShURLFile{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (p *Producer) WriteEvent(event *Event.ShURLFile) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err = p.writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err = p.writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}

func (p *Producer) Close() error {
	// закрываем файл
	return p.file.Close()
}

func (c *Consumer) Close() error {
	// закрываем файл
	return c.file.Close()
}
