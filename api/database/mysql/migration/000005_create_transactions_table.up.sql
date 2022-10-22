CREATE TABLE IF NOT EXISTS `transactions` (
    `id_transaction` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `product_name` VARCHAR(256) NOT NULL,
    `customer_name` VARCHAR(100) NOT NULL,
    `no_transaction` VARCHAR(100) NOT NULL,
    `created_at` TIMESTAMP,
    `updated_at` TIMESTAMP,
    PRIMARY KEY (`id_transaction`)
) ENGINE = InnoDB;