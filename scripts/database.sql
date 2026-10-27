-- 1. Tạo Database mới
CREATE DATABASE IF NOT EXISTS freelance_platform DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE freelance_platform;

-- 2. Bảng Catalogs
CREATE TABLE `complexity` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `complexity_type` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_complexity_type` (`complexity_type`)
);

CREATE TABLE `expected_duration` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `duration_text` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_duration_text` (`duration_text`)
);

CREATE TABLE `payment_type` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `type_name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_payment_type_name` (`type_name`)
);

CREATE TABLE `skill` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `skill_name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_skill_name` (`skill_name`)
);

CREATE TABLE `status_catalog` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `status_name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_status_name` (`status_name`)
);

-- 3. User & Profile
CREATE TABLE `user_account` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `first_name` varchar(100) NOT NULL,
  `last_name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_user_account_email` (`email`)
);

CREATE TABLE `company` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `company_location` varchar(255) DEFAULT NULL,
  `company_name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `freelancer` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `registration_date` date DEFAULT NULL,
  `user_account_id` bigint NOT NULL,
  `location` varchar(255) DEFAULT NULL,
  `overview` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_freelancer_user_account` (`user_account_id`),
  CONSTRAINT `fk_freelancer_user_account` FOREIGN KEY (`user_account_id`) REFERENCES `user_account` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE `hire_manager` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_account_id` bigint NOT NULL,
  `registration_date` date DEFAULT NULL,
  `location` varchar(255) DEFAULT NULL,
  `company_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_hire_manager_user_account` (`user_account_id`),
  KEY `fk_hire_manager_company` (`company_id`),
  CONSTRAINT `fk_hire_manager_company` FOREIGN KEY (`company_id`) REFERENCES `company` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_hire_manager_user_account` FOREIGN KEY (`user_account_id`) REFERENCES `user_account` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

-- 4. Qualification & Test
CREATE TABLE `has_skill` (
  `freelancer_id` bigint NOT NULL,
  `skill_id` bigint NOT NULL,
  PRIMARY KEY (`freelancer_id`,`skill_id`),
  KEY `fk_has_skill_skill` (`skill_id`),
  CONSTRAINT `fk_has_skill_freelancer` FOREIGN KEY (`freelancer_id`) REFERENCES `freelancer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_has_skill_skill` FOREIGN KEY (`skill_id`) REFERENCES `skill` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE `certification` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `certification_name` varchar(255) NOT NULL,
  `freelancer_id` bigint NOT NULL,
  `provider` varchar(255) DEFAULT NULL,
  `description` text DEFAULT NULL,
  `date_earned` date DEFAULT NULL,
  `certification_link` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_certification_freelancer` (`freelancer_id`),
  CONSTRAINT `fk_certification_freelancer` FOREIGN KEY (`freelancer_id`) REFERENCES `freelancer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE `test` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `test_name` varchar(255) NOT NULL,
  `test_link` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_test_name` (`test_name`)
);

CREATE TABLE `test_result` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `freelancer_id` bigint NOT NULL,
  `test_id` bigint NOT NULL,
  `start_time` date NOT NULL,
  `end_time` date DEFAULT NULL,
  `test_result_link` text DEFAULT NULL,
  `score` decimal(5,2) DEFAULT NULL,
  `display_on_profile` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_test_result_freelancer` (`freelancer_id`),
  KEY `idx_test_result_test` (`test_id`),
  CONSTRAINT `fk_test_result_freelancer` FOREIGN KEY (`freelancer_id`) REFERENCES `freelancer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_test_result_test` FOREIGN KEY (`test_id`) REFERENCES `test` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

-- 5. Job & Proposal
CREATE TABLE `job` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `hire_manager_id` bigint NOT NULL,
  `duration_id` bigint NOT NULL,
  `complexity_id` bigint NOT NULL,
  `description` text DEFAULT NULL,
  `main_skill_id` bigint NOT NULL,
  `payment_type_id` bigint NOT NULL,
  `payment_amount` decimal(10,2) DEFAULT NULL,
  `title` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_job_hire_manager` (`hire_manager_id`),
  KEY `idx_job_duration` (`duration_id`),
  KEY `idx_job_complexity` (`complexity_id`),
  KEY `idx_job_main_skill` (`main_skill_id`),
  KEY `idx_job_payment_type` (`payment_type_id`),
  CONSTRAINT `fk_job_complexity` FOREIGN KEY (`complexity_id`) REFERENCES `complexity` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_job_duration` FOREIGN KEY (`duration_id`) REFERENCES `expected_duration` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_job_hire_manager` FOREIGN KEY (`hire_manager_id`) REFERENCES `hire_manager` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_job_main_skill` FOREIGN KEY (`main_skill_id`) REFERENCES `skill` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_job_payment_type` FOREIGN KEY (`payment_type_id`) REFERENCES `payment_type` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE `job_require_skill` (
  `job_id` bigint NOT NULL,
  `skill_id` bigint NOT NULL,
  PRIMARY KEY (`job_id`,`skill_id`),
  KEY `fk_job_require_skill_skill` (`skill_id`),
  CONSTRAINT `fk_job_require_skill_job` FOREIGN KEY (`job_id`) REFERENCES `job` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_job_require_skill_skill` FOREIGN KEY (`skill_id`) REFERENCES `skill` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE `proposal` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `job_id` bigint NOT NULL,
  `freelancer_id` bigint NOT NULL,
  `proposal_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `payment_type_id` bigint NOT NULL,
  `payment_amount` decimal(10,2) NOT NULL,
  `current_proposal_status_id` bigint NOT NULL,
  `client_grade` decimal(2,1) DEFAULT NULL,
  `client_comment` text DEFAULT NULL,
  `freelancer_comment` text DEFAULT NULL,
  `freelancer_grade` decimal(2,1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_proposal_job` (`job_id`),
  KEY `idx_proposal_freelancer` (`freelancer_id`),
  KEY `idx_proposal_payment_type` (`payment_type_id`),
  KEY `idx_proposal_status` (`current_proposal_status_id`),
  CONSTRAINT `fk_proposal_freelancer` FOREIGN KEY (`freelancer_id`) REFERENCES `freelancer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_proposal_job` FOREIGN KEY (`job_id`) REFERENCES `job` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_proposal_payment_type` FOREIGN KEY (`payment_type_id`) REFERENCES `payment_type` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_proposal_status` FOREIGN KEY (`current_proposal_status_id`) REFERENCES `status_catalog` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
);

-- 6. Contract & Work Session
CREATE TABLE `contract` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `proposal_id` bigint NOT NULL,
  `company_id` bigint NOT NULL,
  `freelancer_id` bigint NOT NULL,
  `start_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `end_time` timestamp NULL DEFAULT NULL,
  `payment_type_id` bigint NOT NULL,
  `payment_amount` decimal(10,2) NOT NULL,
  `current_contract_status_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_contract_proposal` (`proposal_id`),
  KEY `idx_contract_company` (`company_id`),
  KEY `idx_contract_freelancer` (`freelancer_id`),
  KEY `idx_contract_payment_type` (`payment_type_id`),
  KEY `idx_contract_status` (`current_contract_status_id`),
  CONSTRAINT `fk_contract_company` FOREIGN KEY (`company_id`) REFERENCES `company` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_contract_freelancer` FOREIGN KEY (`freelancer_id`) REFERENCES `freelancer` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_contract_payment_type` FOREIGN KEY (`payment_type_id`) REFERENCES `payment_type` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_contract_proposal` FOREIGN KEY (`proposal_id`) REFERENCES `proposal` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_contract_status` FOREIGN KEY (`current_contract_status_id`) REFERENCES `status_catalog` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE `work_session` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `contract_id` bigint NOT NULL,
  `start_time` timestamp NULL DEFAULT NULL,
  `end_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_work_session_contract` (`contract_id`),
  CONSTRAINT `fk_work_session_contract` FOREIGN KEY (`contract_id`) REFERENCES `contract` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE `message` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `freelancer_id` bigint DEFAULT NULL,
  `hire_manager_id` bigint DEFAULT NULL,
  `message_time` timestamp NOT NULL,
  `message_text` text COLLATE utf8mb4_general_ci NOT NULL,
  `proposal_id` bigint NOT NULL,
  `proposal_status_catalog_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_message_freelancer` (`freelancer_id`),
  KEY `fk_message_hire_manager` (`hire_manager_id`),
  KEY `fk_message_proposal` (`proposal_id`),
  KEY `fk_message_status_catalog` (`proposal_status_catalog_id`),
  CONSTRAINT `fk_message_freelancer` FOREIGN KEY (`freelancer_id`) REFERENCES `freelancer` (`id`),
  CONSTRAINT `fk_message_hire_manager` FOREIGN KEY (`hire_manager_id`) REFERENCES `hire_manager` (`id`),
  CONSTRAINT `fk_message_proposal` FOREIGN KEY (`proposal_id`) REFERENCES `proposal` (`id`),
  CONSTRAINT `fk_message_status_catalog` FOREIGN KEY (`proposal_status_catalog_id`) REFERENCES `status_catalog` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1254 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- freelancermarket.attachment definition

CREATE TABLE `attachment` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `message_id` bigint NOT NULL,
  `attachment_link` text COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`),
  KEY `fk_attachment_message` (`message_id`),
  CONSTRAINT `fk_attachment_message` FOREIGN KEY (`message_id`) REFERENCES `message` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=188 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- freelancermarket.freelancer_portfolio definition

CREATE TABLE `freelancer_portfolio` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `freelancer_id` bigint NOT NULL,
  `job_id` bigint DEFAULT NULL,
  `project_title` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `project_description` text COLLATE utf8mb4_general_ci,
  `project_url` varchar(500) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `completed_date` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_portfolio_freelancer` (`freelancer_id`),
  KEY `fk_portfolio_job` (`job_id`),
  CONSTRAINT `fk_portfolio_freelancer` FOREIGN KEY (`freelancer_id`) REFERENCES `freelancer` (`id`),
  CONSTRAINT `fk_portfolio_job` FOREIGN KEY (`job_id`) REFERENCES `job` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=677 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;