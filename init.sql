create database if not exists `cat_clerk`;

use `cat_clerk`;

CREATE TABLE `accounts` (
	`id` INT(12) NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(64) NOT NULL COLLATE 'utf8mb4_general_ci',
	`password` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci',
	`salt` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci',
	`email` VARCHAR(64) NOT NULL COLLATE 'utf8mb4_general_ci',
	`dark_theme` TINYINT(1) NOT NULL DEFAULT '1',
	`notifications` TINYINT(1) NOT NULL DEFAULT '0',
	`last_login` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `username` (`username`) USING BTREE,
	UNIQUE INDEX `email` (`email`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;


CREATE TABLE `storages` (
	`id` INT(12) NOT NULL AUTO_INCREMENT,
	`title` VARCHAR(50) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
	`updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `share_requests` (
	`id` INT(12) NOT NULL AUTO_INCREMENT,
	`from_username` VARCHAR(64) NOT NULL COLLATE 'utf8mb4_general_ci',
	`to_username` VARCHAR(64) NOT NULL COLLATE 'utf8mb4_general_ci',
	`share_type` ENUM('storage','shopping_list') NOT NULL COLLATE 'utf8mb4_general_ci',
	`id_request` INT(12) NOT NULL,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `storage_items` (
	`id` INT(12) NOT NULL AUTO_INCREMENT,
	`storage_id` INT(12) NOT NULL,
	`title` VARCHAR(50) NOT NULL COLLATE 'utf8mb4_general_ci',
	`image` VARCHAR(512) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
	`quantity` INT(12) NOT NULL DEFAULT '0',
	`quantity_type` ENUM('grams','kilos','pieces','liters','milliliters') NOT NULL DEFAULT 'pieces' COLLATE 'utf8_general_ci',
	`quantity_threshold` INT(12) NOT NULL DEFAULT '0',
	`expiration_threshold` INT(12) NOT NULL DEFAULT '0',
	`expiration_date` TIMESTAMP NULL DEFAULT NULL,
	`updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `FK_storage_items_storages` (`storage_id`) USING BTREE,
	CONSTRAINT `FK_storage_items_storages` FOREIGN KEY (`storage_id`) REFERENCES `cat_clerk`.`storages` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `shopping_lists` (
	`id` INT(12) NOT NULL AUTO_INCREMENT,
	`title` VARCHAR(50) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
	`updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `shopping_list_items` (
	`id` INT(12) NOT NULL AUTO_INCREMENT,
	`shopping_list_id` INT(12) NOT NULL,
	`title` VARCHAR(50) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
	`quantity` INT(12) NOT NULL DEFAULT '1',
	`updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `FK_shopping_list_items_shopping_lists` (`shopping_list_id`) USING BTREE,
	CONSTRAINT `FK_shopping_list_items_shopping_lists` FOREIGN KEY (`shopping_list_id`) REFERENCES `cat_clerk`.`shopping_lists` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `settings` (
	`id` INT(12) NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(64) NOT NULL COLLATE 'utf8mb4_general_ci',
	`dark_theme` TINYINT(1) NOT NULL DEFAULT '1',
	`notifications` TINYINT(1) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `username` (`username`) USING BTREE,
	CONSTRAINT `FK__accounts` FOREIGN KEY (`username`) REFERENCES `cat_clerk`.`accounts` (`username`) ON UPDATE CASCADE ON DELETE CASCADE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `foods` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(50) NOT NULL COLLATE 'utf8mb4_general_ci',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `name` (`name`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `account_storage_binder` (
	`id` INT(12) NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(64) NOT NULL COLLATE 'utf8mb4_general_ci',
	`storage_id` INT(12) NOT NULL,
	`owner` TINYINT(1) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `username_storage_id` (`username`, `storage_id`) USING BTREE,
	INDEX `FK_account_storage_binder_storages` (`storage_id`) USING BTREE,
	INDEX `FK_account_storage_binder_accounts` (`username`) USING BTREE,
	CONSTRAINT `FK_account_storage_binder_accounts` FOREIGN KEY (`username`) REFERENCES `cat_clerk`.`accounts` (`username`) ON UPDATE NO ACTION ON DELETE NO ACTION,
	CONSTRAINT `FK_account_storage_binder_storages` FOREIGN KEY (`storage_id`) REFERENCES `cat_clerk`.`storages` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `account_shopping_list_binder` (
	`id` INT(12) NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(64) NOT NULL COLLATE 'utf8mb4_general_ci',
	`shopping_list_id` INT(12) NOT NULL,
	`owner` TINYINT(1) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `username_shopping_list_id` (`username`, `shopping_list_id`) USING BTREE,
	INDEX `FK_account_shopping_list_binder_shopping_lists` (`shopping_list_id`) USING BTREE,
	CONSTRAINT `FK_account_shopping_list_binder_accounts` FOREIGN KEY (`username`) REFERENCES `cat_clerk`.`accounts` (`username`) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT `FK_account_shopping_list_binder_shopping_lists` FOREIGN KEY (`shopping_list_id`) REFERENCES `cat_clerk`.`shopping_lists` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;