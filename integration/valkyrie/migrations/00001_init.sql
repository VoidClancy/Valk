-- +goose Up
CREATE TABLE `User` (
  `id` text NOT NULL,
  `email` text NOT NULL,
  `phoneNum` text NOT NULL,
  `role` text NOT NULL DEFAULT ('student'),
  `referredById` text NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `User_referredById_fkey` FOREIGN KEY (`referredById`) REFERENCES `User` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `User_role_check` CHECK ("role" IN ('ADMIN', 'student', 'TEACHER'))
);
CREATE UNIQUE INDEX `User_email_key` ON `User` (`email`);
CREATE UNIQUE INDEX `User_email_phoneNum_key` ON `User` (`email`, `phoneNum`);
CREATE TABLE `Profile` (
  `id` text NOT NULL,
  `bio` text NULL,
  `userId` text NOT NULL,
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
  `text` text NOT NULL,
  `postId` text NOT NULL,
  `authorId` text NOT NULL,
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

-- +goose Down
PRAGMA foreign_keys = off;
DROP TABLE `User`;
DROP TABLE `Profile`;
DROP TABLE `Post`;
DROP TABLE `Comment`;
DROP TABLE `Category`;
DROP TABLE `CategoryToPost`;
PRAGMA foreign_keys = on;
