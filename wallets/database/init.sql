CREATE TABLE `wallets`
(
    id   bigint auto_increment,
    ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    iin varchar(255) NOT NULL,
    accountno varchar(255) NOT NULL UNIQUE,
    amount int DEFAULT 0,
    PRIMARY KEY (`id`)
);

INSERT INTO `wallets`(`accountno`, `iin`) VALUES('KZT5004100100', '0');