CREATE DATABASE sntest;
USE sntest;

CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(100) NOT NULL,
  `phash` varbinary(32) NOT NULL,
  `name` varchar(100),
  `surname` varchar(100),
  `birthdate` datetime,
  `gender` tinyint(1) DEFAULT 0,
  `city` varchar(200),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `posts` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user` bigint(20) NOT NULL,
  `header` varchar(100),
  `updated` datetime,
  `text` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `d_interests` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(100),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `interests` (
  `user` bigint(20) NOT NULL,
  `interest` bigint(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `friends` (
  `user` bigint(20) NOT NULL,
  `friend` bigint(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
