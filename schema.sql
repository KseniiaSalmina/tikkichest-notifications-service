CREATE TABLE IF NOT EXISTS notifications (
    "profile_id" INT PRIMARY KEY NOT NULL,
    "telegram_username" TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS profile_id_idx_notifications_idx ON notifications(profile_id);
CREATE INDEX IF NOT EXISTS telegram_username_idx_notifications_idx ON notifications USING HASH(telegram_username);