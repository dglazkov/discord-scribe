DROP TABLE `messages`;
CREATE TABLE `messages` (
  `id` bigint(8) NOT NULL,
  `channel_id` bigint(8) NOT NULL,
  `guild_id` bigint(8) DEFAULT NULL,
  `author_id` bigint(8) DEFAULT NULL,
  `content` text NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
