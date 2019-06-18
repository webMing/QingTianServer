/*
 Navicat MariaDB Data Transfer

 Source Server         : Cont
 Source Server Type    : MariaDB
 Source Server Version : 50560
 Source Host           : localhost:3306
 Source Schema         : qingtian

 Target Server Type    : MariaDB
 Target Server Version : 50560
 File Encoding         : 65001

 Date: 18/06/2019 22:22:14
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `user_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT 'xiaoqingtian',
  `user_id` int(100) NOT NULL AUTO_INCREMENT,
  `passwd` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone_num` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `union_id` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `user_client` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `user_client_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=20190618 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
