CREATE TABLE IF NOT EXISTS `products` (
    `id_product` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(10) NOT NULL UNIQUE,
    `name` VARCHAR(100) NOT NULL,
    `description` VARCHAR(100) NOT NULL,
    `price` INTEGER(11) NOT NULL,
    `image` VARCHAR(50),
    `created_at` TIMESTAMP,
    `updated_at` TIMESTAMP,
    PRIMARY KEY (`id_product`)
) ENGINE = InnoDB;