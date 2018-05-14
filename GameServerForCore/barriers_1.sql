/*
Navicat MySQL Data Transfer

Source Server         : 192.168.0.189
Source Server Version : 50617
Source Host           : 192.168.0.189:3306
Source Database       : barriers_1

Target Server Type    : MYSQL
Target Server Version : 50617
File Encoding         : 65001

Date: 2018-04-23 17:06:36
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for playerdata
-- ----------------------------
DROP TABLE IF EXISTS `playerdata`;
CREATE TABLE `playerdata` (
  `pid` bigint(19) NOT NULL,
  `userid` int(11) NOT NULL,
  `serverid` int(11) NOT NULL,
  `pname` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `picon` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `piconborder` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `pexp` int(11) DEFAULT NULL,
  `pdiamond` int(11) DEFAULT NULL,
  `pgold` int(11) DEFAULT NULL,
  `pspirit` int(11) DEFAULT NULL,
  `pmedal` int(11) DEFAULT NULL,
  `phonor` int(11) DEFAULT NULL,
  `pcreatetime` datetime(6) DEFAULT NULL,
  `plastlogintime` datetime(6) DEFAULT NULL,
  `prandseed` int(11) DEFAULT NULL,
  `pguildid` int(11) DEFAULT NULL,
  `plastspiritsddtime` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`pid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Records of playerdata
-- ----------------------------

-- ----------------------------
-- Table structure for playerequip
-- ----------------------------
DROP TABLE IF EXISTS `playerequip`;
CREATE TABLE `playerequip` (
  `pid` bigint(19) NOT NULL,
  `sid` int(11) NOT NULL,
  `equipid` int(11) DEFAULT NULL,
  `eexp` int(11) DEFAULT NULL,
  `erand` int(11) DEFAULT NULL,
  `eawakenrand` int(11) DEFAULT NULL,
  `eawakencount` int(11) DEFAULT NULL,
  `ecreateticks` datetime(6) DEFAULT NULL,
  `eowner` int(11) DEFAULT NULL,
  `eislock` bit(1) DEFAULT NULL,
  `estate` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`pid`,`sid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Records of playerequip
-- ----------------------------

-- ----------------------------
-- Table structure for playerhero
-- ----------------------------
DROP TABLE IF EXISTS `playerhero`;
CREATE TABLE `playerhero` (
  `pid` bigint(19) NOT NULL,
  `heroid` int(11) NOT NULL,
  `hexp` int(11) DEFAULT NULL,
  `hquality` int(11) DEFAULT NULL,
  `hawaken` int(11) DEFAULT NULL,
  `hequip1` int(11) DEFAULT NULL,
  `hequip2` int(11) DEFAULT NULL,
  `hequip3` int(11) DEFAULT NULL,
  `hequip4` int(11) DEFAULT NULL,
  `hequip5` int(11) DEFAULT NULL,
  `hequip6` int(11) DEFAULT NULL,
  `hskillrock1` int(11) DEFAULT NULL,
  `hskillrock2` int(11) DEFAULT NULL,
  `hskillrock3` int(11) DEFAULT NULL,
  `hskillrock4` int(11) DEFAULT NULL,
  `hskillrock5` int(11) DEFAULT NULL,
  `hskillrock6` int(11) DEFAULT NULL,
  `hskillrock7` int(11) DEFAULT NULL,
  `hskillrock8` int(11) DEFAULT NULL,
  `hskin` int(11) DEFAULT NULL,
  PRIMARY KEY (`pid`,`heroid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Records of playerhero
-- ----------------------------

-- ----------------------------
-- Table structure for playeritem
-- ----------------------------
DROP TABLE IF EXISTS `playeritem`;
CREATE TABLE `playeritem` (
  `pid` bigint(19) NOT NULL,
  `itemid` int(11) NOT NULL,
  `num` int(11) DEFAULT NULL,
  PRIMARY KEY (`pid`,`itemid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Records of playeritem
-- ----------------------------

-- ----------------------------
-- Table structure for playerrock
-- ----------------------------
DROP TABLE IF EXISTS `playerrock`;
CREATE TABLE `playerrock` (
  `pid` bigint(19) NOT NULL,
  `rockid` int(11) NOT NULL,
  `num` int(11) DEFAULT NULL,
  PRIMARY KEY (`pid`,`rockid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Records of playerrock
-- ----------------------------
