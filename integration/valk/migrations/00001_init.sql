-- +goose Up
CREATE TABLE `User` (
  `id` text NOT NULL,
  `email` text NOT NULL,
  `phoneNum` text NOT NULL,
  `password` text NULL,
  `role` text NOT NULL DEFAULT ('student'),
  `roleOptional` text NULL,
  `loginCount` integer NOT NULL DEFAULT (0),
  `referredById` text NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `User_referredById_fkey` FOREIGN KEY (`referredById`) REFERENCES `User` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `User_role_check` CHECK ("role" IN ('ADMIN', 'student', 'TEACHER')),
  CONSTRAINT `User_roleOptional_check` CHECK ("roleOptional" IN ('ADMIN', 'student', 'TEACHER'))
);
CREATE UNIQUE INDEX `User_email_key` ON `User` (`email`);
CREATE UNIQUE INDEX `User_phoneNum_key` ON `User` (`phoneNum`);
CREATE UNIQUE INDEX `emailPhone` ON `User` (`email`, `phoneNum`);
CREATE TABLE `Profile` (
  `id` text NOT NULL,
  `bio` text NULL,
  `userId` text NOT NULL,
  `createdAt` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`id`),
  CONSTRAINT `Profile_userId_fkey` FOREIGN KEY (`userId`) REFERENCES `User` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE UNIQUE INDEX `Profile_userId_key` ON `Profile` (`userId`);
CREATE TABLE `Post` (
  `id` text NOT NULL,
  `title` text NOT NULL,
  `content` text NULL,
  `published` integer NOT NULL DEFAULT (FALSE),
  `authorId` text NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `Post_authorId_fkey` FOREIGN KEY (`authorId`) REFERENCES `User` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE `Comment` (
  `id` text NOT NULL,
  `textify` integer NOT NULL,
  `dummy3` text NOT NULL,
  `dummy1` integer NOT NULL,
  `dummy2` text NOT NULL,
  `postId` text NOT NULL,
  `authorId` text NOT NULL,
  `meta` blob NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `Comment_postId_fkey` FOREIGN KEY (`postId`) REFERENCES `Post` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `Comment_authorId_fkey` FOREIGN KEY (`authorId`) REFERENCES `User` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE `Category` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL
);
CREATE UNIQUE INDEX `Category_name_key` ON `Category` (`name`);
CREATE TABLE `CategoryToPost` (
  `postId` text NOT NULL,
  `categoryId` integer NOT NULL,
  PRIMARY KEY (`postId`, `categoryId`),
  CONSTRAINT `CategoryToPost_postId_fkey` FOREIGN KEY (`postId`) REFERENCES `Post` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `CategoryToPost_categoryId_fkey` FOREIGN KEY (`categoryId`) REFERENCES `Category` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE `DefaultsTest` (
  `uuid4` text NOT NULL,
  `uuid7` text NOT NULL,
  `uuidNoArgs` text NOT NULL,
  `cuid1` text NOT NULL,
  `cuid2` text NOT NULL,
  `cuidNoArgs` text NOT NULL,
  `ulid` text NOT NULL,
  `nanoid` text NOT NULL,
  `now` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`uuid4`)
);
CREATE TABLE `AllFieldsSoFar` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `stringReq` text NOT NULL,
  `stringOpt` text NULL,
  `stringDefault` text NOT NULL DEFAULT ('default'),
  `stringVarchar` text NOT NULL,
  `stringChar` text NOT NULL,
  `bitVal` text NOT NULL,
  `varBitVal` text NOT NULL,
  `inetVal` text NOT NULL,
  `xmlVal` text NOT NULL,
  `cuidDefault` text NOT NULL,
  `cuid1Default` text NOT NULL,
  `cuid2Default` text NOT NULL,
  `uuidDefault` text NOT NULL,
  `uuid4Default` text NOT NULL,
  `uuid7Default` text NOT NULL,
  `ulidDefault` text NOT NULL,
  `nanoidDefault` text NOT NULL,
  `uuidDb` text NOT NULL,
  `intReq` integer NOT NULL,
  `intOpt` integer NULL,
  `intDefault` integer NOT NULL DEFAULT (42),
  `integerVal` integer NOT NULL,
  `smallInt` integer NOT NULL,
  `tinyInt` integer NOT NULL,
  `oidVal` integer NOT NULL,
  `bigIntReq` integer NOT NULL,
  `bigIntOpt` integer NULL,
  `floatReq` real NOT NULL,
  `floatOpt` real NULL,
  `realVal` real NOT NULL,
  `decimalReq` real NOT NULL,
  `decimalOpt` real NULL,
  `decimalPrecise` real NOT NULL,
  `moneyVal` real NOT NULL,
  `boolReq` integer NOT NULL,
  `boolOpt` integer NULL,
  `boolDefault` integer NOT NULL DEFAULT (FALSE),
  `dateTimeReq` timestamp NOT NULL,
  `dateTimeOpt` timestamp NULL,
  `dateTimeDefault` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updatedAt` timestamp NOT NULL,
  `dateTimeTz` timestamp NOT NULL,
  `timestampVal` timestamp NOT NULL,
  `timeVal` timestamp NOT NULL,
  `timetzVal` timestamp NOT NULL,
  `jsonReq` blob NOT NULL,
  `jsonOpt` blob NULL,
  `jsonVal` blob NOT NULL,
  `bytesReq` blob NOT NULL,
  `bytesOpt` blob NULL
);

-- +goose Down
PRAGMA foreign_keys = off;
DROP TABLE `User`;
DROP TABLE `Profile`;
DROP TABLE `Post`;
DROP TABLE `Comment`;
DROP TABLE `Category`;
DROP TABLE `CategoryToPost`;
DROP TABLE `DefaultsTest`;
DROP TABLE `AllFieldsSoFar`;
PRAGMA foreign_keys = on;
