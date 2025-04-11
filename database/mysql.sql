SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE IF NOT EXISTS wxtoken CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
use wxtoken;

CREATE TABLE mp_access_token (
	id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  app_id BIGINT UNSIGNED NOT NULL,
  mp_id VARCHAR(64) NOT NULL,
  access_token VARCHAR(512) NOT NULL,
  deadline BIGINT UNSIGNED NOT NULL,
  expires_in INT,
  created_at BIGINT NOT NULL,
	updated_at BIGINT NOT NULL,
	deleted_at DATETIME,
	`version`  BIGINT,
  PRIMARY KEY (id),
  INDEX `idx_app_id` (app_id),
  INDEX `idx_mp_id` (mp_id),
  UNIQUE INDEX `idx_app_id_mp_id` (app_id, mp_id)
) AUTO_INCREMENT = 1000 COMMENT = "AccessToken表";

CREATE TABLE platform_app (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(32) NOT NULL,
    `type` INT NOT NULL,
    token VARCHAR(32) NOT NULL,
    encoding_aes_key VARCHAR(48) NOT NULL,
    encoding_type INT NOT NULL,
    app_id VARCHAR(48) NOT NULL,
    app_secret VARCHAR(64) NOT NULL,
    `status` INT NOT NULL DEFAULT 0,
    introduction VARCHAR(32),
    pic_url VARCHAR(256),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at DATETIME,
    `version` BIGINT,
    PRIMARY KEY (id),
    INDEX `idx_app_id` (app_id)
)  AUTO_INCREMENT=1001;

SET FOREIGN_KEY_CHECKS = 1;
