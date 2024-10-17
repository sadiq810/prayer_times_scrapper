/*
SQLyog Community v13.3.0 (64 bit)
MySQL - 10.4.32-MariaDB : Database - golang_prayer_times
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
/*Table structure for table `cities` */

DROP TABLE IF EXISTS `cities`;

CREATE TABLE `cities` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `country_id` int(11) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `image` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `f_id` int(11) DEFAULT NULL,
  `f_status` tinyint(1) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `country_id` (`country_id`),
  KEY `title` (`title`),
  KEY `f_id` (`f_id`),
  KEY `f_status` (`f_status`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

/*Data for the table `cities` */

/*Table structure for table `countries` */

DROP TABLE IF EXISTS `countries`;

CREATE TABLE `countries` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL,
  `image` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `f_id` int(11) DEFAULT NULL COMMENT 'used id of remote data',
  `f_status` tinyint(1) DEFAULT 0 COMMENT 'show status of the remote data fetching',
  PRIMARY KEY (`id`),
  KEY `f_id` (`f_id`),
  KEY `f_status` (`f_status`),
  KEY `title` (`title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

/*Data for the table `countries` */

/*Table structure for table `masjid_prayer_times` */

DROP TABLE IF EXISTS `masjid_prayer_times`;

CREATE TABLE `masjid_prayer_times` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `masjid_id` int(11) NOT NULL,
  `prayer_name_id` int(11) NOT NULL,
  `date` date NOT NULL,
  `adhan_time` varchar(255) NOT NULL,
  `iqamah_time` varchar(255) DEFAULT NULL,
  `month` int(11) DEFAULT NULL,
  `day` int(11) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `masjid_prayer_times_masjid_id_index` (`masjid_id`),
  KEY `prayer_name_id_idx` (`prayer_name_id`),
  KEY `month` (`month`),
  KEY `day` (`day`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*Data for the table `masjid_prayer_times` */

/*Table structure for table `masjids` */

DROP TABLE IF EXISTS `masjids`;

CREATE TABLE `masjids` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `country_id` int(11) DEFAULT NULL,
  `city_id` int(11) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `lat` varchar(255) DEFAULT NULL,
  `lng` varchar(255) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `image` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `f_id` bigint(21) DEFAULT NULL,
  `f_guid` varchar(255) DEFAULT NULL,
  `f_status` tinyint(1) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `country_id` (`country_id`),
  KEY `city_id` (`city_id`),
  KEY `lat` (`lat`),
  KEY `lng` (`lng`),
  KEY `f_id` (`f_id`),
  KEY `f_guid` (`f_guid`),
  KEY `f_status` (`f_status`),
  KEY `title` (`title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*Data for the table `masjids` */

/*Table structure for table `prayer_names` */

DROP TABLE IF EXISTS `prayer_names`;

CREATE TABLE `prayer_names` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL,
  `sort_order` int(11) DEFAULT 0,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `sort_order` (`sort_order`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*Data for the table `prayer_names` */

insert  into `prayer_names`(`id`,`title`,`sort_order`,`created_at`,`updated_at`) values 
(1,'Fajr',1,'2024-01-06 01:05:02','2024-01-06 01:05:12'),
(2,'Dhuhr',3,'2024-01-06 01:05:04','2024-01-06 01:05:14'),
(3,'Asr',4,'2024-01-06 01:05:06','2024-01-06 01:05:16'),
(4,'Maghrib',5,'2024-01-06 01:05:08','2024-01-06 01:05:17'),
(5,'Isha',6,'2024-01-06 01:05:10','2024-01-06 01:05:19'),
(6,'Friday prayer',7,'2024-01-06 01:05:44','2024-01-06 01:05:47'),
(9,'shouruq',2,'2024-10-04 14:58:07','2024-10-04 14:58:09');

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
