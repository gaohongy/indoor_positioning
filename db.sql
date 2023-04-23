-- 建库
CREATE DATABASE `indoor_positioning` CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_general_ci';

-- 建场所表
CREATE TABLE `indoor_positioning`.`place`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '场所id',
  `place_address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '场所详细地址',
  `longitude` decimal(10,6) NOT NULL COMMENT '经度',
  `latitude` decimal(10,6) NOT NULL COMMENT '纬度',
  `createdate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedate` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`)
);

-- 建用户表
CREATE TABLE `indoor_positioning`.`user`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
  `salt` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '盐',
  `pwdhash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户密码',
  `usertype` tinyint NOT NULL COMMENT '用户类型',
  `place_id` int NULL COMMENT '场所id',
  `createdate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedate` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  FOREIGN KEY (`place_id`) REFERENCES `indoor_positioning`.`place` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT
);


-- 建网格点表
CREATE TABLE `indoor_positioning`.`gridpoint`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '网格点id',
  `coordinate_x` decimal(10,6) NOT NULL COMMENT 'x坐标',
  `coordinate_y` decimal(10,6) NOT NULL COMMENT 'y坐标',
  `coordinate_z` int NOT NULL COMMENT 'z坐标(楼层)',
  `place_id` int NULL COMMENT '场所id',
  `createdate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedate` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  FOREIGN KEY (`place_id`) REFERENCES `indoor_positioning`.`place` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
);


-- 建路径点表
CREATE TABLE `indoor_positioning`.`pathpoint`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '路径点id',
  `user_id` int NOT NULL COMMENT '用户id',
  `grid_point_id` int NOT NULL COMMENT '网格点id',
  `place_id` int NULL COMMENT '场所id',
  `createdate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedate` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES `indoor_positioning`.`user` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  FOREIGN KEY (`grid_point_id`) REFERENCES `indoor_positioning`.`grid_point` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
  FOREIGN KEY (`place_id`) REFERENCES `indoor_positioning`.`place` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
);


-- 建AP表
CREATE TABLE `indoor_positioning`.`ap`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'AP_id',
  `ssid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '网络名称',
  `bssid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'AP_MAC地址',
  `grid_point_id` int NOT NULL COMMENT '网格点id',
  `place_id` int NOT NULL COMMENT '场所id',
  `createdate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedate` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  FOREIGN KEY (`grid_point_id`) REFERENCES `indoor_positioning`.`grid_point` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  FOREIGN KEY (`place_id`) REFERENCES `indoor_positioning`.`place` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
);

-- 建参考点表
CREATE TABLE `indoor_positioning`.`referencepoint`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '参考点id',
  `grid_point_id` int NOT NULL COMMENT '网格点id',
  `place_id` int NOT NULL COMMENT '场所id',
  `createdate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedate` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  FOREIGN KEY (`grid_point_id`) REFERENCES `indoor_positioning`.`grid_point` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  FOREIGN KEY (`place_id`) REFERENCES `indoor_positioning`.`place` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
);

-- 建rss表
CREATE TABLE `indoor_positioning`.`rss`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '信号强度id',
  `rss` decimal(10,6) NOT NULL COMMENT '信号强度',
  `reference_point_id` int NOT NULL COMMENT '参考点id',
  `ap_id` int NOT NULL COMMENT 'AP_id',
  `createdate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedate` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  FOREIGN KEY (`reference_point_id`) REFERENCES `indoor_positioning`.`reference_point` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  FOREIGN KEY (`ap_id`) REFERENCES `indoor_positioning`.`ap` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
);