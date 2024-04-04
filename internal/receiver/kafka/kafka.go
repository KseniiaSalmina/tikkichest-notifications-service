package kafka

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/IBM/sarama"

	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/config"
	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/notifier"
)

type ConsumerManager struct {
	manager       sarama.ConsumerGroup
	topic         string
	consumer      Consumer
	finishClosing sync.WaitGroup
}

type Consumer struct {
	messageCh chan notifier.Notification
	closeCtx  context.Context
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	close(c.messageCh)
	return nil
}
func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Println("message channel closed")
				return nil
			}

			user, err := strconv.Atoi(string(msg.Key))
			if err != nil {
				return fmt.Errorf("failed to get profile id: %w", err)
			}

			formatMsg := notifier.Notification{
				User:    user,
				Message: msg.Value,
			}
			c.messageCh <- formatMsg

			sess.MarkMessage(msg, "notifier")

		case <-c.closeCtx.Done():
			return nil

		case <-sess.Context().Done():
			return nil
		}
	}
}

func NewConsumerManager(cfg config.Kafka) (*ConsumerManager, error) {
	consumerGroup, err := sarama.NewConsumerGroup([]string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}, cfg.ConsumerGroupID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer manager: %w", err)
	}

	return &ConsumerManager{
		manager:       consumerGroup,
		topic:         cfg.Topic,
		finishClosing: sync.WaitGroup{},
	}, nil
}

func (cm *ConsumerManager) Run(ctx context.Context) <-chan notifier.Notification {
	notifications := make(chan notifier.Notification)

	go func(ch chan notifier.Notification) {
		cm.finishClosing.Add(1)
		defer cm.finishClosing.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := cm.manager.Consume(ctx, []string{cm.topic}, &Consumer{messageCh: ch, closeCtx: ctx}); err != nil {
					log.Println(err) //TODO logger
				}
			}
		}
	}(notifications)

	return notifications
}

func (cm *ConsumerManager) Shutdown() error {
	err := cm.manager.Close()
	cm.finishClosing.Wait()

	return err
}
