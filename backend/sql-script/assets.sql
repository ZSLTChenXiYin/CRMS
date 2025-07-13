CREATE TABLE `crms`.`assets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `type` enum('server') NOT NULL,
  `name` varchar(191) NOT NULL,
  `data` json NOT NULL,
  `owner_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_assets_owner_id` (`owner_id`),
  KEY `idx_assets_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;