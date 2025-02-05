package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"weather-notification/internal/domain/entity"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = "notifications"
	queueName    = "notifications.send"
	retryQueue   = "notifications.retry"
	dlqQueue     = "notifications.dlq"
)

type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQService(url string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar canal: %w", err)
	}

	service := &RabbitMQService{
		conn:    conn,
		channel: ch,
	}

	if err := service.setup(); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *RabbitMQService) setup() error {
	err := s.channel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("erro ao declarar exchange: %w", err)
	}

	queues := []string{queueName, retryQueue, dlqQueue}
	for _, q := range queues {
		_, err = s.channel.QueueDeclare(
			q,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("erro ao declarar fila %s: %w", q, err)
		}

		err = s.channel.QueueBind(
			q,
			q,
			exchangeName,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("erro ao bind fila %s: %w", q, err)
		}
	}

	return nil
}

func (s *RabbitMQService) PublishNotification(ctx context.Context, notification *entity.Notification) error {
	log.Printf("Serializando notificação %s", notification.ID)
	data, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("erro ao serializar notificação: %w", err)
	}
	log.Printf("Publicando notificação no RabbitMQ para fila %s", queueName)

	err = s.channel.PublishWithContext(ctx,
		exchangeName,
		queueName,
		true,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			ContentType:  "application/json",
			Body:         data,
		},
	)

	if err != nil {
		log.Printf("Erro ao publicar no RabbitMQ: %v", err)
		return fmt.Errorf("erro ao publicar notificação: %w", err)
	}

	log.Printf("Notificação publicada com sucesso no RabbitMQ")
	return nil
}

func (s *RabbitMQService) ConsumeNotifications(ctx context.Context, handler func(*entity.Notification) error) error {
	msgs, err := s.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	retryMsgs, err := s.channel.Consume(retryQueue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go s.processMessages(msgs, handler)
	go s.processMessages(retryMsgs, handler)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-forever:
		return nil
	}
}

func (s *RabbitMQService) processMessages(msgs <-chan amqp.Delivery, handler func(*entity.Notification) error) {
	for msg := range msgs {
		var notification entity.Notification
		if err := json.Unmarshal(msg.Body, &notification); err != nil {
			msg.Reject(false)
			continue
		}

		if time.Now().Before(notification.ScheduledFor) {
			msg.Reject(true)
			continue
		}

		err := handler(&notification)
		if err != nil {
			s.handleError(msg)
			msg.Reject(false)
			continue
		}

		msg.Ack(false)
	}
}

func (s *RabbitMQService) handleError(msg amqp.Delivery) {
	retryCount := 0
	if val, ok := msg.Headers["x-retry-count"]; ok {
		switch v := val.(type) {
		case int:
			retryCount = v
		case int32:
			retryCount = int(v)
		case int64:
			retryCount = int(v)
		case float64:
			retryCount = int(v)
		default:
			fmt.Printf("Tipo inesperado para x-retry-count: %T\n", v)
		}
	}

	if count, ok := msg.Headers["x-retry-count"].(float64); ok {
		retryCount = int(count)
	}

	if retryCount >= 3 {
		s.publishToDLQ(msg.Body)
		return
	}

	s.channel.Publish(
		exchangeName,
		retryQueue,
		false,
		false,
		amqp.Publishing{
			Headers: amqp.Table{
				"x-retry-count": retryCount + 1,
			},
			Body:         msg.Body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (s *RabbitMQService) publishToDLQ(body []byte) {
	s.channel.Publish(
		exchangeName,
		dlqQueue,
		false,
		false,
		amqp.Publishing{
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (s *RabbitMQService) Close() error {
	if err := s.channel.Close(); err != nil {
		return fmt.Errorf("erro ao fechar canal: %w", err)
	}
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("erro ao fechar conexão: %w", err)
	}
	return nil
}
