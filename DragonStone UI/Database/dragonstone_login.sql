-- phpMyAdmin SQL Dump
-- version 4.6.5.2
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Oct 25, 2017 at 09:37 PM
-- Server version: 10.1.21-MariaDB
-- PHP Version: 5.6.30

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `dragonstone_login`
--

-- --------------------------------------------------------

--
-- Table structure for table `user`
--

CREATE TABLE `user` (
  `UserId` int(11) NOT NULL,
  `UserName` varchar(255) NOT NULL,
  `UserPassword` varchar(255) NOT NULL,
  `UserEmail` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='Mashup User';

--
-- Dumping data for table `user`
--

INSERT INTO `user` (`UserId`, `UserName`, `UserPassword`, `UserEmail`) VALUES
(1, 'admin', 'admin', 'admin@gmail.com'),
(2, 'renu', 'renu', 'renu@gmail.com'),
(3, 'vivek', 'vivek123', 'vivek@live.com'),
(4, 'Smitha V', 'Smitha123', 'smithav17@gmail.com'),
(5, 'karan', '12BF1a04c6', 'karan@kukka.com'),
(6, 'google', 'google', 'google@gmail.com'),
(7, 'Chaitanya Kademane', 'chaitu123', 'chait.2605@gmail.com'),
(8, 'Chaitanya Kademane', 'chaitu123', 'chait.2605@gmail.com'),
(9, 'Chaitanya Kademane', 'chaitu123', 'chait.2605@gmail.com'),
(10, 'Chaitanya Kademane', 'chaitu123', 'chait.2605@gmail.com'),
(11, 'Chaitu', 'chaitu123', 'chaitu123@gmail.com'),
(12, 'cai', 'cai123', 'cai@gmail.com'),
(13, 'Aaditya Deowanshi', 'qwerty', 'aaditya.deowanshi@gmail.com'),
(14, 'Aaditya Deowanshi', 'qwerty', 'aaditya.deowanshi@gmail.com'),
(15, 'Aaditya Deowanshi', 'qwerty', 'aaditya.deowanshi@gmail.com'),
(16, 'Aaditya . Deowanshi', '', 'aaditya.deowanshi@sjsu.edu'),
(17, 'Smitha Venkatesh', '', 'smitha.venkatesh@sjsu.edu'),
(18, 'Karthik Nair', '', 'karthikreads@gmail.com'),
(19, 'spartan deals', '', 'spartandeals6@gmail.com'),
(20, 'renu parameswaran', '', 'renu.parameswaran@gmail.com'),
(21, 'Arun Ram', '', 'y.arunram@gmail.com'),
(22, 'Jyothi Parameswaran', '', 'rupaparameswar@gmail.com');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`UserId`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `user`
--
ALTER TABLE `user`
  MODIFY `UserId` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=23;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
