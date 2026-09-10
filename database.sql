SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

CREATE SCHEMA IF NOT EXISTS `mydb` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `mydb`;

CREATE TABLE IF NOT EXISTS `mydb`.`user_account` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `user_name` VARCHAR(100) NOT NULL UNIQUE,
    `password_hash` VARCHAR(255) NOT NULL,
    `first_name` VARCHAR(100) NOT NULL,
    `last_name` VARCHAR(100) NOT NULL,
    PRIMARY KEY (`id`)
)  ENGINE=INNODB;


CREATE TABLE IF NOT EXISTS `mydb`.`freelancer` (
    `id` BIGINT UNSIGNED NOT NULL UNIQUE,
    `registration_date` DATE,
    `user_account_id` BIGINT UNSIGNED NOT NULL UNIQUE,
    `location` VARCHAR(255) NULL,
    `overview` TEXT NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_account_id`)
        REFERENCES user_account (`id`)
)  ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `mydb`.`certification` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `certification_name` VARCHAR(255) NOT NULL,
    `freelancer_id` BIGINT UNSIGNED NULL UNIQUE,
    `provider` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `date_earned` DATE,
    `certinfication_link` TEXT,
    FOREIGN KEY (`freelancer_id`)
        REFERENCES freelancer (`id`)
)  ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `mydb`.`skill` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `skill_name` VARCHAR(255) NOT NULL UNIQUE
)  ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `mydb`.`has_skill` (
    `freelancer_id` BIGINT UNSIGNED NOT NULL,
    `skill_id` BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`freelancer_id` , `skill_id`),
    FOREIGN KEY (`freelancer_id`)
        REFERENCES freelancer (`id`),
    FOREIGN KEY (`skill_id`)
        REFERENCES skill (`id`)
)  ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `mydb`.`test` (
    `test_id` BIGINT UNSIGNED NOT NULL,
    `test_name` VARCHAR(255) NOT NULL UNIQUE
    `test_link` VARCHAR(500),
    PRIMARY KEY(test_id)
)

CREATE TABLE IF NOT EXISTS `mydb`.`test_result` (
    `test_result_id` BIGINT UNSIGNED NOT NULL,
    `freelancer_id` BIGINT NOT NULL, --FK
    `test_id`BIGINT NOT NULL, --fk
    `start_time` DATE NOT NULL,
    `end_time` DATE NULL,
    `test_result_link` TEXT NULL,
    `score` DECIMAL(5,2) NULL,
    `display_on_profile` BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY(test_result_id),
    FOREIGN KEY(`freelancer_id`) REFERENCES freelancer(`id`),
    FOREIGN KEY(`test_id`) REFERENCES test(`id`)
)

CREATE TABLE IF NOT EXISTS `mydb`.`company` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `company_location` VARCHAR(255),
    `company_name` VARCHAR(255) NOT NULL
)  ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `mydb`.`hire_manager` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `user_account_id` BIGINT UNSIGNED NOT NULL,
    `registration_date` DATE NOT NULL,
    `location` VARCHAR(255),
    `company_id` BIGINT UNSIGNED NULL,
    FOREIGN KEY (`user_account_id`)
        REFERENCES user_account (`id`),
    FOREIGN KEY (`company_id`)
        REFERENCES company (`id`)
)  ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `mydb`.`payment_type` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `type_name` VARCHAR(100) UNIQUE
);

CREATE TABLE IF NOT EXISTS `mydb`.`proposal_status_catalog` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `status_name` VARCHAR(255) NOT NULL UNIQUE 
);

CREATE TABLE IF NOT EXISTS `mydb`.`expected_duration` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `duration_text` VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS `mydb`.`complexity` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `complexity_type` VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS `mydb`.`job` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `hire_manager_id` BIGINT UNSIGNED NOT NULL,
    `duration_id` BIGINT UNSIGNED NOT NULL,
    `complexity_id` BIGINT UNSIGNED NOT NULL,
    `description` TEXT NOT NULL,
    `main_skill_id` BIGINT UNSIGNED NOT NULL,
    `payment_type_id` BIGINT UNSIGNED NOT NULL,
    `payment_amount` DECIMAL(12,2) NOT NULL,
    FOREIGN KEY (`hire_manager_id`)
        REFERENCES hire_manager (`id`),
    FOREIGN KEY (`duration_id`)
        REFERENCES expected_duration (`id`),
    FOREIGN KEY (`complexity_id`)
        REFERENCES complexity (`id`),
    FOREIGN KEY (`main_skill_id`)
        REFERENCES skill (`id`),
    FOREIGN KEY (`payment_type_id`)
        REFERENCES payment_type (`id`)
);

CREATE TABLE IF NOT EXISTS `mydb`.`job_require_skill` (
    `job_id` BIGINT UNSIGNED NOT NULL,
    `skill_id` BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`job_id` , `skill_id`),
    FOREIGN KEY (`job_id`)
        REFERENCES job (`id`),
    FOREIGN KEY (`skill_id`)
        REFERENCES skill (`id`)
);

CREATE TABLE IF NOT EXISTS `mydb`.`proposal` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `job_id` BIGINT UNSIGNED NOT NULL,
    `freelancer_id` BIGINT UNSIGNED NOT NULL,
    `proposal_time` TIMESTAMP,
    `payment_type_id` BIGINT UNSIGNED NOT NULL,
    `payment_amount` DECIMAL(8 ,2) NOT NULL,
    `current_proposal_status_id` BIGINT UNSIGNED NOT NULL,
    `client_grade` DECIMAL(1 , 1 ) CHECK (client_grade <= 5.0
        AND client_grade >= 0.0),
    `client_comment` TEXT,
    `freelancer_comment` TEXT,
    `freelancer_grade` DECIMAL(1 , 1 ) CHECK (freelancer_grade <= 5.0
        AND freelancer_grade >= 0.0),
    FOREIGN KEY (`job_id`)
        REFERENCES job (`id`),
    FOREIGN KEY (`freelancer_id`)
        REFERENCES freelancer (`id`),
    FOREIGN KEY (`payment_type_id`)
        REFERENCES payment_type (`id`),
    FOREIGN KEY (`current_proposal_status_id`)
        REFERENCES proposal_status_catalog (`id`)
);


CREATE TABLE IF NOT EXISTS `mydb`.`contract` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `proposal_id` BIGINT UNSIGNED NOT NULL UNIQUE,
    `company_id` BIGINT UNSIGNED NOT NULL,
    `freelancer_id` BIGINT UNSIGNED NOT NULL,
    `start_time` TIMESTAMP NOT NULL,
    `end_time` TIMESTAMP NULL,
    `payment_type_id` BIGINT UNSIGNED NOT NULL,
    `payment_amount` DECIMAL(8 , 2 ) NOT NULL,
    FOREIGN KEY (`proposal_id`)
        REFERENCES proposal (`id`),
    FOREIGN KEY (`company_id`)
        REFERENCES company (`id`),
    FOREIGN KEY (`freelancer_id`)
        REFERENCES freelancer (`id`),
    FOREIGN KEY (`payment_type_id`)
        REFERENCES payment_type (`id`)
)  ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `mydb`.`work_session` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `contract_id` BIGINT UNSIGNED NOT NULL,
    `start_time` TIMESTAMP,
    `end_time` TIMESTAMP NULL,
    FOREIGN KEY (`contract_id`)
        REFERENCES contract (`id`)
)  ENGINE=INNODB;




CREATE TABLE IF NOT EXISTS `mydb`.`message` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `freelancer_id` BIGINT UNSIGNED NULL,
    `hire_manager_id` BIGINT UNSIGNED NULL,
    `message_time` TIMESTAMP NOT NULL,
    `message_text` TEXT NOT NULL,
    `proposal_id` BIGINT UNSIGNED NOT NULL,
    `proposal_status_catalog_id` BIGINT UNSIGNED NULL,
    FOREIGN KEY (`freelancer_id`)
        REFERENCES freelancer (`id`),
    FOREIGN KEY (`hire_manager_id`)
        REFERENCES hire_manager (`id`),
    FOREIGN KEY (`proposal_id`)
        REFERENCES proposal (`id`),
    FOREIGN KEY (`proposal_status_catalog_id`)
        REFERENCES proposal_status_catalog (`id`)
);



CREATE TABLE IF NOT EXISTS `mydb`.`attachment` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `message_id` BIGINT UNSIGNED NOT NULL,
    `attachment_link` TEXT,
    FOREIGN KEY (`message_id`)
        REFERENCES message (`id`)
);


