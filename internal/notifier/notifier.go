package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type Notification struct {
	User    int
	Message []byte
}

type Sender interface {
	SendNotification(username string, event Event) error
}

type Receiver interface {
	Run(ctx context.Context) <-chan Notification
}

type Storage interface {
	GetUsername(ctx context.Context, id int) (string, error)
}

type Notifier struct {
	sender        Sender
	receiver      Receiver
	storage       Storage
	finishClosing sync.WaitGroup
}

func NewNotifier(sender Sender, receiver Receiver, storage Storage) *Notifier {
	return &Notifier{
		sender:        sender,
		receiver:      receiver,
		storage:       storage,
		finishClosing: sync.WaitGroup{},
	}
}

func (n *Notifier) Run(ctx context.Context) {
	notificationCh := n.receiver.Run(ctx)

	go n.worker(ctx, notificationCh)
	n.finishClosing.Add(1)
}

func (n *Notifier) worker(ctx context.Context, ch <-chan Notification) {
	defer n.finishClosing.Done()

	for {
		select {
		case notification := <-ch:
			username, err := n.getUsernameByID(ctx, notification.User)
			if err != nil {
				log.Println(err) //TODO логгер
				continue
			}

			var event Event
			if err := json.Unmarshal(notification.Message, &event); err != nil {
				log.Println(err) //TODO логгер
				continue
			}

			if err := validateEvent(event); err != nil {
				log.Println(err) //TODO логгер
				continue
			}

			if err := n.sender.SendNotification(username, event); err != nil {
				log.Println(err) //TODO логгер
			}

		case <-ctx.Done():
			return
		}
	}
}

func (n *Notifier) getUsernameByID(ctx context.Context, id int) (string, error) {
	username, err := n.storage.GetUsername(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to find username by id %d: %w", id, err)
	}

	return username, nil
}

func (n *Notifier) Shutdown() {
	n.finishClosing.Wait()
}
