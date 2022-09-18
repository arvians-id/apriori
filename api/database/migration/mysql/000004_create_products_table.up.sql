CREATE TABLE IF NOT EXISTS `products` (
    `id_product` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(10) NOT NULL UNIQUE,
    `name` VARCHAR(100) NOT NULL,
    `description` TEXT,
    `price` INTEGER NOT NULL,
    `category` VARCHAR(100) NOT NULL,
    `is_empty` TINYINT(1) NOT NULL DEFAULT 0,
    `mass` INTEGER NOT NULL,
    `image` TEXT,
    `created_at` TIMESTAMP,
    `updated_at` TIMESTAMP,
    PRIMARY KEY (`id_product`)
) ENGINE = InnoDB;