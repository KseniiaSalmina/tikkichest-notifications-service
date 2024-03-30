package notifier

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Notification struct {
	User    int
	Message string
}

type Sender interface {
	SendNotification(username, message string) error
}

type Updater interface {
	Run(ctx context.Context) ([]<-chan Notification, error)
}

type Storage interface {
	GetUsername(ctx context.Context, id int) (string, error)
}

type Notifier struct {
	sender        Sender
	updater       Updater
	storage       Storage
	finishClosing *sync.WaitGroup
}

func NewNotifier(sender Sender, updater Updater, storage Storage) *Notifier {
	return &Notifier{
		sender:  sender,
		updater: updater,
		storage: storage,
	}
}

func (n *Notifier) Run(ctx context.Context) error {
	notificationChans, err := n.updater.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to start notifier: %w", err)
	}

	for _, ch := range notificationChans {
		go n.worker(ctx, ch)
		n.finishClosing.Add(1)
	}

	return nil
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

			if err := n.sender.SendNotification(username, notification.Message); err != nil {
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
