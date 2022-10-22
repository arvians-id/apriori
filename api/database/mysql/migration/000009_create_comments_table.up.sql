CREATE TABLE IF NOT EXISTS `comments` (
    `id_comment` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_order_id` BIGINT(20) UNSIGNED NOT NULL,
    `product_code` VARCHAR(10) NOT NULL,
    `description` TEXT,
    `tag` VARCHAR(200),
    `rating` INTEGER NOT NULL,
    `created_at` TIMESTAMP,
    PRIMARY KEY (`id_comment`),
    FOREIGN KEY (`user_order_id`) REFERENCES user_orders(`id_order`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB;