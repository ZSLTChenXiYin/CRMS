-- 带外键约束的改进版
CREATE TABLE `crms`.`users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `email` varchar(191) NOT NULL,
  `password_hash` varchar(191) NOT NULL,
  `expired_at` datetime(3) DEFAULT NULL,
  `last_login_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB ROW_FORMAT=DYNAMIC DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `crms`.`assets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `type` enum('server') NOT NULL,
  `name` varchar(191) NOT NULL,
  `data` json NOT NULL,
  `owner_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_assets_owner_id` (`owner_id`),
  KEY `idx_assets_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_asset_owner` FOREIGN KEY (`owner_id`) 
    REFERENCES `crms`.`users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB ROW_FORMAT=DYNAMIC DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `crms`.`user_asset_mappings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `asset_id` bigint unsigned NOT NULL,
  `permission` enum('use','execute') NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_asset` (`user_id`,`asset_id`),
  KEY `idx_user_asset_mappings_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_mapping_user` FOREIGN KEY (`user_id`) 
    REFERENCES `crms`.`users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_mapping_asset` FOREIGN KEY (`asset_id`) 
    REFERENCES `crms`.`assets` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB ROW_FORMAT=DYNAMIC DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `crms`.`operation_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `resource_type` enum('user','asset') NOT NULL,
  `resource_id` bigint unsigned NOT NULL,
  `action` enum('create','update','delete') NOT NULL,
  `action_details` text NOT NULL,
  `additional_information` json DEFAULT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_operation_logs_user_id` (`user_id`),
  KEY `idx_operation_logs_resource_id` (`resource_id`),
  CONSTRAINT `fk_log_user` FOREIGN KEY (`user_id`) 
    REFERENCES `crms`.`users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB ROW_FORMAT=DYNAMIC DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;