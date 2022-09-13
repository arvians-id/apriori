CREATE TABLE IF NOT EXISTS notifications (
    id_notification SERIAL,
    user_id INTEGER NOT NULL,
    title VARCHAR(100) NOT NULL ,
    description TEXT,
    url VARCHAR(200),
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP,
    PRIMARY KEY (id_notification),
    FOREIGN KEY (user_id) REFERENCES users(id_user) ON DELETE CASCADE ON UPDATE CASCADE
)