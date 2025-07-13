CREATE TABLE `crms`.`operation_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `resource_type` enum('user','asset') NOT NULL,
  `resource_id` bigint unsigned NOT NULL,
  `action` enum('create','update','delete') NOT NULL,
  `action_details` text NOT NULL,
  `additional_information` json DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_operation_logs_user_id` (`user_id`),
  KEY `idx_operation_logs_resource_id` (`resource_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;