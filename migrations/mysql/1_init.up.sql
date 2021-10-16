CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(100) NOT NULL,
  `phash` varbinary(32) NOT NULL,
  `name` varchar(100),
  `surname` varchar(100),
  `birthdate` datetime,
  `gender` tinyint(1) DEFAULT 0,
  `city` varchar(200),
  CONSTRAINT users_uc_username UNIQUE(username),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `posts` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user` bigint(20) NOT NULL,
  `header` varchar(100),
  `updated` datetime,
  `text` text,
  UNIQUE KEY (`header`,`user`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `d_interests` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(100),
  CONSTRAINT d_interests_uc_name UNIQUE(name),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `interests` (
  `user` bigint(20) NOT NULL,
  `interest` bigint(20) NOT NULL,
  CONSTRAINT interests_uc_user_interest UNIQUE(user, interest)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `friends` (
  `user` bigint(20) NOT NULL,
  `friend` bigint(20) NOT NULL,
  CONSTRAINT friends_uc_user_friend UNIQUE(user, friend)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `subscribers` (
  `user` bigint(20) NOT NULL,
  `subscriber` bigint(20) NOT NULL,
  CONSTRAINT subscribers_uc_user_subscriber UNIQUE(user, subscriber)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
