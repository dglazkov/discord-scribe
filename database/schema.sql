DROP TABLE IF EXISTS `channels`;

CREATE TABLE `channels` (
  `id` bigint(8) NOT NULL,
  `has_beginning` boolean DEFAULT 0, -- contains messages from the very beginning
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `messages`;

CREATE TABLE `messages` (
  `id` bigint(8) NOT NULL,
  `channel_id` bigint(8) NOT NULL,
  `guild_id` bigint(8) DEFAULT NULL,
  `author_id` bigint(8) DEFAULT NULL,
  `content` text NOT NULL,
  `timestamp` timestamp NOT NULL,
  `reaction_count` int DEFAULT 0,
  `reaction_types` int DEFAULT 0,
  `type` int DEFAULT 0,
  -- TODO: consider making channel_id a foreign key
  PRIMARY KEY (`id`),
  INDEX (`timestamp`, `channel_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
