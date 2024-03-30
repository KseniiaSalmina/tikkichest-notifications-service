package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/config"
)

type DB struct {
	db *pgxpool.Pool
}

func NewDB(ctx context.Context, cfg config.Postgres) (*DB, error) {
	connstr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := pgxpool.New(ctx, connstr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		db: db,
	}, nil
}

func (db *DB) SaveUsername(ctx context.Context, id int, username string) error {
	if _, err := db.db.Exec(ctx, `INSERT INTO notifications (profile_id, telegram_username) VALUES ($1, $2)`, id, username); err != nil {
		return fmt.Errorf("failed to save username: %w", err)
	}

	return nil
}

func (db *DB) DeleteUsername(ctx context.Context, id int) error {
	if _, err := db.db.Exec(ctx, `DELETE FROM notifications WHERE profile_id = $1`, id); err != nil {
		return fmt.Errorf("failed to delete username: %w", err)
	}

	return nil
}

func (db *DB) ChangeUsername(ctx context.Context, id int, username string) error {
	if _, err := db.db.Exec(ctx, `UPDATE notifications SET telegram_username = $1 WHERE profile_id = $2`, username, id); err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	return nil
}

func (db *DB) GetUsername(ctx context.Context, id int) (string, error) {
	var username pgtype.Text
	if err := db.db.QueryRow(ctx, `SELECT telegram_username FROM notifications WHERE profile_id=$1`, id).Scan(&username); err != nil {
		return "", fmt.Errorf("failed to get username: %w", err)
	}

	return username.String, nil
}

func (db *DB) Close() {
	db.db.Close()
}
