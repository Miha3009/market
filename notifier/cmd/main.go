package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/miha3009/market/notifier/internal/service"
	"github.com/miha3009/market/notifier/pkg/config"
	"github.com/segmentio/kafka-go"
)

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	cfg, err := config.ReadConfig("config.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)},
		Topic:     cfg.Kafka.Topic,
		Partition: cfg.Kafka.Partition,
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffsetAt(context.TODO(), time.Now())

	emailService := service.NewMailService(logger, cfg.Server)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		emailService.SendMail(string(m.Key), "Market", string(m.Value))
		logger.Println("Mail sended to " + string(m.Key))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
