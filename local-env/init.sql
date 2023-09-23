DROP DATABASE IF EXISTS app;
CREATE DATABASE app;
USE app;
CREATE TABLE `customer`
(
    `id`           int          NOT NULL AUTO_INCREMENT,
    `name`         varchar(255) NOT NULL,
    `email`        varchar(255) NOT NULL,
    `password`     varchar(255) NOT NULL,
    `phone_number` varchar(15) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `account`
(
    `id`          int            NOT NULL AUTO_INCREMENT,
    `customer_id` int            NOT NULL,
    `balance`     decimal(10, 2) NOT NULL DEFAULT '0.00',
    `deleted`     int            NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY           `customer_id` (`customer_id`,`deleted`),
    CONSTRAINT `account_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `transaction`
(
    `id`               int            NOT NULL AUTO_INCREMENT,
    `account_id`       int            NOT NULL,
    `transaction_type` enum('Deposit','Withdrawal','Transfer') NOT NULL,
    `amount`           decimal(10, 2) NOT NULL,
    `timestamp`        datetime DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY                `account_id` (`account_id`),
    CONSTRAINT `transaction_ibfk_1` FOREIGN KEY (`account_id`) REFERENCES `account` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

