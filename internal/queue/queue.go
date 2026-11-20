package queue

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/vviveksharma/lorebound/internal/logger"
	"github.com/vviveksharma/lorebound/internal/utils"
)

// Config holds Kafka client settings from the environment.
type Config struct {
	Brokers []string
	Topic   string
	GroupID string
}

// LoadConfig reads KAFKA_* variables from the environment.
func LoadConfig() Config {
	brokers := strings.Split(utils.MustGetenv("KAFKA_BROKERS"), ",")
	for i := range brokers {
		brokers[i] = strings.TrimSpace(brokers[i])
	}
	return Config{
		Brokers: brokers,
		Topic:   utils.MustGetenv("KAFKA_INGEST_TOPIC"),
		GroupID: utils.Getenv("KAFKA_CONSUMER_GROUP", "lorebound-workers"),
	}
}

// Producer publishes messages to the configured ingest topic.
type Producer struct {
	writer *kafka.Writer
	topic  string
}

// Connect verifies the cluster, ensures the ingest topic exists, and returns a producer.
func Connect() (*Producer, error) {
	cfg := LoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := ping(ctx, cfg.Brokers); err != nil {
		return nil, fmt.Errorf("kafka ping: %w", err)
	}
	if err := ensureTopic(ctx, cfg.Brokers, cfg.Topic); err != nil {
		return nil, fmt.Errorf("kafka ensure topic %q: %w", cfg.Topic, err)
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}

	logger.Info("kafka connected", "brokers", cfg.Brokers, "topic", cfg.Topic)
	return &Producer{writer: writer, topic: cfg.Topic}, nil
}

func ping(ctx context.Context, brokers []string) error {
	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	_, err = conn.Brokers()
	return err
}

func ensureTopic(ctx context.Context, brokers []string, topic string) error {
	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	controllerConn, err := kafka.DialContext(ctx, "tcp", controllerAddr)
	if err != nil {
		return err
	}
	defer func() { _ = controllerConn.Close() }()

	return controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
}

// Publish sends a single message to the ingest topic.
func (p *Producer) Publish(ctx context.Context, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}

// Close releases the underlying writer.
func (p *Producer) Close() error {
	return p.writer.Close()
}

// Topic returns the ingest topic name.
func (p *Producer) Topic() string {
	return p.topic
}

// NewConsumer builds a reader for the ingest topic using KAFKA_CONSUMER_GROUP.
func NewConsumer(cfg Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
		Topic:    cfg.Topic,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
}
