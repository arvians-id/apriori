CREATE TABLE IF NOT EXISTS `apriories` (
    `id_apriori` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(10) NOT NULL,
    `item` VARCHAR(256) NOT NULL,
    `discount` DECIMAL(6,2) NOT NULL,
    `support` DECIMAL(6,2) NOT NULL,
    `confidence` DECIMAL(6,2) NOT NULL,
    `range_date` VARCHAR(50) NOT NULL,
    `is_active` TINYINT(1) NOT NULL DEFAULT 0,
    `description` TEXT,
    `image` TEXT,
    `created_at` TIMESTAMP,
    PRIMARY KEY (`id_apriori`)
) ENGINE = InnoDB;