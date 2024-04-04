package app

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/api"
	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/config"
	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/notifier"
	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/receiver/kafka"
	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/sender/telegram"
	postgresql "github.com/KseniiaSalmina/tikkichest-notifications-service/internal/storage/postgres"
)

type Application struct {
	cfg          config.Application
	db           *postgresql.DB
	sender       *telegram.Bot
	receiver     *kafka.ConsumerManager
	notifier     *notifier.Notifier
	server       *api.Server
	closeCtx     context.Context
	closeCtxFunc context.CancelFunc
}

func NewApplication(cfg config.Application) (*Application, error) {
	app := Application{
		cfg: cfg,
	}

	if err := app.bootstrap(); err != nil {
		return nil, err
	}

	app.readyToShutdown()

	return &app, nil
}

func (a *Application) bootstrap() error {
	//init dependencies
	if err := a.initDatabase(); err != nil {
		return fmt.Errorf("failed to bootstrap application: %w", err)
	}

	//init services
	a.initSender()
	if err := a.initReceiver(); err != nil {
		return fmt.Errorf("failed to bootstrap application: %w", err)
	}

	//init controllers
	a.initNotifier()
	a.initServer()

	return nil
}

func (a *Application) initDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := postgresql.NewDB(ctx, a.cfg.Postgres)
	if err != nil {
		log.Println(err) // TODO: logger
		return fmt.Errorf("failed to init database: %w", err)
	}

	a.db = db
	log.Println("successful connection to database") // TODO: logger
	return nil
}

func (a *Application) initSender() {
	a.sender = telegram.NewBot(a.cfg.Telegram)
}

func (a *Application) initReceiver() error {
	receiver, err := kafka.NewConsumerManager(a.cfg.Kafka)
	if err != nil {
		return fmt.Errorf("failed to init receiver: %w", err)
	}

	a.receiver = receiver
	return nil
}

func (a *Application) initNotifier() {
	a.notifier = notifier.NewNotifier(a.sender, a.receiver, a.db)
}

func (a *Application) initServer() {
	a.server = api.NewServer(a.cfg.Server, a.db)
}

func (a *Application) Run() {
	defer a.stop()

	a.server.Run()
	a.notifier.Run(a.closeCtx)

	<-a.closeCtx.Done()
	a.closeCtxFunc()
}

func (a *Application) stop() {
	if err := a.server.Shutdown(); err != nil {
		log.Printf("incorrect closing of server: %s", err.Error()) // TODO: logger
	} else {
		log.Print("server closed") // TODO: logger
	}

	a.receiver.Shutdown()
	a.notifier.Shutdown()

	a.db.Close()
	log.Print("database closed") // TODO: logger
}

func (a *Application) readyToShutdown() {
	ctx, closeCtx := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	a.closeCtx = ctx
	a.closeCtxFunc = closeCtx
}
