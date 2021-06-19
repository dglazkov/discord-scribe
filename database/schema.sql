DROP TABLE IF EXISTS `channels`;

CREATE TABLE `channels` (
  `id` bigint(8) NOT NULL,
  `is_fully_read` boolean,
  `earliest_read_message` bigint(8) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `messages`;

CREATE TABLE `messages` (
  `id` bigint(8) NOT NULL,
  `channel_id` bigint(8) NOT NULL,
  `guild_id` bigint(8) DEFAULT NULL,
  `author_id` bigint(8) DEFAULT NULL,
  `content` text NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  -- TODO: consider making channel_id a foreign key
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
