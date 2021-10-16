ALTER TABLE posts ADD COLUMN `created` datetime;
ALTER TABLE posts ADD COLUMN `deleted` tinyint(1) DEFAULT 0;
ALTER TABLE users ADD COLUMN `rebuild_feed_flag` tinyint(1) DEFAULT 0;