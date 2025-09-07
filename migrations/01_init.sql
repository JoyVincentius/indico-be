CREATE DATABASE `indico` /*!40100 DEFAULT CHARACTER SET latin1 */;

-- indico.job_records definition

CREATE TABLE `job_records` (
  `id` varchar(191) NOT NULL,
  `status` longtext,
  `progress` bigint(20) DEFAULT NULL,
  `processed` bigint(20) DEFAULT NULL,
  `total` bigint(20) DEFAULT NULL,
  `result_path` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `cancelled` tinyint(1) DEFAULT NULL,
  `cancel_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- indico.jobs definition

CREATE TABLE `jobs` (
  `id` varchar(255) NOT NULL,
  `status` varchar(50) DEFAULT 'QUEUED',
  `progress` int(11) DEFAULT '0',
  `processed` int(11) DEFAULT '0',
  `total` int(11) DEFAULT '0',
  `result_path` varchar(255) DEFAULT NULL,
  `cancelled` tinyint(1) DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- indico.orders definition

CREATE TABLE `orders` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `quantity` bigint(20) DEFAULT NULL,
  `buyer_id` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_product` (`product_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

-- indico.products definition

CREATE TABLE `products` (
  `id` bigint(20) NOT NULL,
  `stock` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- indico.settlements definition

CREATE TABLE `settlements` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint(20) unsigned DEFAULT NULL,
  `date` datetime(3) DEFAULT NULL,
  `gross_cents` bigint(20) DEFAULT NULL,
  `fee_cents` bigint(20) DEFAULT NULL,
  `net_cents` bigint(20) DEFAULT NULL,
  `txn_count` bigint(20) DEFAULT NULL,
  `generated_at` datetime(3) DEFAULT NULL,
  `run_id` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_merchant_date` (`merchant_id`,`date`),
  KEY `idx_merchant_date` (`merchant_id`,`date`)
) ENGINE=InnoDB AUTO_INCREMENT=258973 DEFAULT CHARSET=latin1;

-- indico.transactions definition

CREATE TABLE `transactions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint(20) unsigned DEFAULT NULL,
  `amount_cents` bigint(20) DEFAULT NULL,
  `fee_cents` bigint(20) DEFAULT NULL,
  `status` longtext,
  `paid_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_merchant_date` (`merchant_id`,`paid_at`)
) ENGINE=InnoDB AUTO_INCREMENT=1000001 DEFAULT CHARSET=latin1;

-- seed

INSERT INTO indico.products
(id, stock)
VALUES(1, 99);