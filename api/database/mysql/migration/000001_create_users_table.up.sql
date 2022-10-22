CREATE TABLE IF NOT EXISTS `users` (
    `id_user` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `role` INT(11) DEFAULT 2 NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `email` VARCHAR(100) NOT NULL UNIQUE,
    `address` VARCHAR(100) NOT NULL,
    `phone` VARCHAR(100) NOT NULL,
    `password` VARCHAR(256) NOT NULL,
    `created_at` TIMESTAMP,
    `updated_at` TIMESTAMP,
    PRIMARY KEY (`id_user`)
) ENGINE = InnoDB;