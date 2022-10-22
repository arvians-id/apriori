CREATE TABLE IF NOT EXISTS notifications (
    `id_notification` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT(20) UNSIGNED NOT NULL,
    `title` VARCHAR(100) NOT NULL,
    `description` TEXT,
    `url` VARCHAR(200),
    `is_read` TINYINT(1) NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP,
    PRIMARY KEY (`id_notification`),
    FOREIGN KEY (`user_id`) REFERENCES users(`id_user`) ON DELETE RESTRICT ON UPDATE CASCADE
)