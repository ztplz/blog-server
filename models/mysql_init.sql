# 初始化数据库

# 如果数据库 blog 不存在就建立一个叫 blog 的数据库，字符集 utf8， 大小写不敏感
CREATE DATABASE IF NOT EXISTS blog CHARACTER SET utf8 COLLATE utf8_general_ci;



# 如果表 admin 不存在就建立一个叫 admin 的表
CREATE TABLE IF NOT EXISTS `admin` (
  `id`         INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `admin_id`   varchar(255)     NOT NULL DEFAULT '',
  `password`   varchar(255)     NOT NULL DEFAULT '',
  `admin_name` varchar(255)     NOT NULL DEFAULT '',
  -- `email` varchar(255) NOT NULL DEFAULT '',
  `image` varchar(255) NOT NULL DEFAULT '',
  `last_login_at` varchar(255) NOT NULL DEFAULT '',
  -- `token` varchar(255) NOT NULL DEFAULT '',
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