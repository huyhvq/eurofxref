CREATE TABLE IF NOT EXISTS `rates`
(
    `id`         integer PRIMARY KEY AUTO_INCREMENT,
    `currency`   varchar(3)     NOT NULL,
    `rate`       decimal(10, 5) NOT NULL,
    `created_at` date           NOT NULL
);