CREATE DATABASE entrytask DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user` (
                        `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                        `username` varchar(100) NOT NULL,
                        `password` varchar(100) NOT NULL,
                        `nickname` varchar(50) DEFAULT NULL,
                        `picfile` varchar(150) DEFAULT NULL,
                        PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE INDEX username ON user(username);