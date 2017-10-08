# 初始化数据库

# 如果数据库 blog 不存在就建立一个叫 blog 的数据库，字符集 utf8， 大小写不敏感
CREATE DATABASE IF NOT EXISTS blog CHARACTER SET utf8 COLLATE utf8_general_ci;



# 如果表 admin 不存在就建立一个叫 admin 的表
CREATE TABLE IF NOT EXISTS `admin` (
  `id`         INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `admin_id`   varchar(255)     NOT NULL DEFAULT '',
  `password`   varchar(255)     NOT NULL DEFAULT '',
  `admin_name` varchar(255)     NOT NULL DEFAULT '',
  `image` varchar(255) NOT NULL DEFAULT '',
  `last_login_at` varchar(255) NOT NULL DEFAULT '',
  `ip` varchar(255) NOT NULL DEFAULT '',
  primary key (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# 如果表 article 不存在就建立一个叫 article 的表
CREATE TABLE IF NOT EXISTS `article` (
  `id`                  INT(11) UNSIGNED        NOT NULL AUTO_INCREMENT,
  `create_at`           datetime                NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at`           datetime                NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `visit_count`         INT(11) UNSIGNED        NOT NULL DEFAULT 0,
  `reply_count`         INT(11) UNSIGNED        NOT NULL DEFAULT 0,
  `article_title`       varchar(255)            NOT NULL DEFAULT '',
  `article_previewtext` text                    NOT NULL,
  `article_content`     text                    NOT NULL,
  `top`                 TINYINT(1)              NOT NULL DEFAULT 0,
  `category`            INT(11)  UNSIGNED       NOT NULL DEFAULT 0,
  `tag_list`            varchar(255)            NOT NULL DEFAULT '',
  primary key (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# 如果表 category_list 不存在就奖励一个叫 category_list 的表
CREATE TABLE IF NOT EXISTS `category_list` (
  `id`         INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `category`   varchar(255)     NOT NULL DEFAULT '',
  primary key (id)        
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# 如果表 tags 不存在就奖励一个叫 tags 的表
CREATE TABLE IF NOT EXISTS `tags` (
  `id`              INT(11) UNSIGNED    NOT NULL AUTO_INCREMENT,
  `color`           varchar(255)        NOT NULL DEFAULT '',
  `tag_title`       varchar(255)        NOT NULL DEFAULT '',
  -- `article_id`      INT(11) UNSIGNED    NOT NULL DEFAULT 0,
  primary key (id)        
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# 如果表 visitor_count 不存在就奖励一个叫 visit_count 的表
CREATE TABLE IF NOT EXISTS `visitor_count` (
  `id`              INT(11) UNSIGNED    NOT NULL AUTO_INCREMENT,
  `date`            varchar(255)        NOT NULL DEFAULT '',
  `count`           INT(11) UNSIGNED        NOT NULL DEFAULT 0,
  -- `article_id`      INT(11) UNSIGNED    NOT NULL DEFAULT 0,
  primary key (id)        
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




# INSERT INTO `admin` SET admin_id='admin', password='$2a$10$uIBcwGWH.hwsilRGd34HHOTOmzbiWUD4buHziGK59TDxJNXSv26cW', admin_name='admin';

# CREATE TABLE `userinfo` (
#   `uid` INT(10) NOT NULL AUTO_INCREMENT,
#   `username` VARCHAR(64) NULL DEFAULT NULL,
#   `departname` VARCHAR(64) NULL DEFAULT NULL,
#   `created` DATE NULL DEFAULT NULL,
#   PRIMARY KEY (`uid`)
# );
#
# CREATE TABLE `userdetail` (
#   `uid` INT(10) NOT NULL DEFAULT '0',
#   `intro` TEXT NULL,
#   `profile` TEXT NULL,
#   PRIMARY KEY (`uid`)
# )