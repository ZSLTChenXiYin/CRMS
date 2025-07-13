CREATE TABLE `crms`.`users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `email` varchar(191) NOT NULL,
  `password_hash` varchar(191) NOT NULL,
  `expired_at` datetime(3) DEFAULT NULL,
  `last_login_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `crms`.`user_asset_mappings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `asset_id` bigint unsigned NOT NULL,
  `permission` enum('use','execute') NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_asset` (`user_id`,`asset_id`),
  KEY `idx_user_asset_mappings_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;