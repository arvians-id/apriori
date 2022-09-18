CREATE TABLE IF NOT EXISTS `password_resets` (
    `email` VARCHAR(100) NOT NULL UNIQUE,
    `token` VARCHAR(256) NOT NULL,
    `expired` INTEGER NOT NULL
) ENGINE = InnoDB;