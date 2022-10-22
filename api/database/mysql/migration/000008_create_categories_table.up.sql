CREATE TABLE IF NOT EXISTS `categories` (
    `id_category` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(20) NOT NULL,
    `created_at` TIMESTAMP,
    `updated_at` TIMESTAMP,
    PRIMARY KEY (`id_category`)
) ENGINE = InnoDB;