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
	manager       sarama.Consumer
	topic         string
	consumers     []sarama.PartitionConsumer
	finishClosing *sync.WaitGroup
}

func NewConsumerManager(cfg config.Kafka) (*ConsumerManager, error) {
	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer manager: %w", err)
	}

	return &ConsumerManager{
		manager: consumer,
		topic:   cfg.Topic,
	}, nil
}

func (cm *ConsumerManager) Run(ctx context.Context) ([]<-chan notifier.Notification, error) {
	if err := cm.createConsumers(); err != nil {
		return nil, fmt.Errorf("failed to start consumer manager: %w", err)
	}

	messageChannels := make([]<-chan notifier.Notification, len(cm.consumers))
	for _, consumer := range cm.consumers {
		messageChannels = append(messageChannels, cm.startListening(ctx, consumer.Messages()))
	}

	return messageChannels, nil
}

func (cm *ConsumerManager) startListening(ctx context.Context, receive <-chan *sarama.ConsumerMessage) <-chan notifier.Notification {
	send := make(chan notifier.Notification)

	go func(send chan notifier.Notification, receive <-chan *sarama.ConsumerMessage) {
		cm.finishClosing.Add(1)
		for {
			select {
			case msg := <-receive:
				user, err := strconv.Atoi(string(msg.Key))
				if err != nil {
					log.Println(err) //TODO логгер
					continue
				}
				message := string(msg.Value)

				formatMsg := notifier.Notification{
					User:    user,
					Message: message,
				}
				send <- formatMsg

			case <-ctx.Done():
				close(send)
				cm.finishClosing.Done()
				return
			}
		}
	}(send, receive)

	return send
}

func (cm *ConsumerManager) createConsumers() error {
	partitions, err := cm.manager.Partitions(cm.topic)
	if err != nil {
		return fmt.Errorf("failed to get partitions IDs: %w", err)
	}

	cm.consumers = make([]sarama.PartitionConsumer, 0, len(partitions))
	for _, id := range partitions {
		consumer, err := cm.newConsumer(id)
		if err != nil {
			return fmt.Errorf("failed to create consumer group: %w", err)
		}
		cm.consumers = append(cm.consumers, consumer)
	}

	return nil
}

func (cm *ConsumerManager) newConsumer(partition int32) (sarama.PartitionConsumer, error) {
	consumer, err := cm.manager.ConsumePartition(cm.topic, partition, sarama.OffsetOldest)
	if err != nil {
		return nil, fmt.Errorf("failed to create partition consumer: %w", err)
	}

	return consumer, nil
}

func (cm *ConsumerManager) Shutdown() []error {
	cm.finishClosing.Wait()

	errs := make([]error, 0, len(cm.consumers))
	for _, consumer := range cm.consumers {
		if err := consumer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("closing partition consumer error: %w/n", err))
		}
	}

	if err := cm.manager.Close(); err != nil {
		errs = append(errs, fmt.Errorf("closing consumer manager error: %w", err))
	}

	return errs
}
