-- +goose Up
PRAGMA foreign_keys = off;
CREATE TABLE `new_User` (
  `id` text NOT NULL,
  `email` text NOT NULL,
  `phoneNum` text NOT NULL,
  `role` text NOT NULL DEFAULT ('student'),
  `referredById` text NULL,
  `sponsorId` text NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `User_referredById_fkey` FOREIGN KEY (`referredById`) REFERENCES `User` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `User_sponsorId_fkey` FOREIGN KEY (`sponsorId`) REFERENCES `User` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `User_role_check` CHECK ("role" IN ('ADMIN', 'student', 'TEACHER'))
);
INSERT INTO `new_User` (`id`, `email`, `phoneNum`, `role`, `referredById`) SELECT `id`, `email`, `phoneNum`, `role`, `referredById` FROM `User`;
DROP TABLE `User`;
ALTER TABLE `new_User` RENAME TO `User`;
CREATE UNIQUE INDEX `User_email_key` ON `User` (`email`);
CREATE UNIQUE INDEX `User_email_phoneNum_key` ON `User` (`email`, `phoneNum`);
PRAGMA foreign_keys = on;

-- +goose Down
PRAGMA foreign_keys = off;
CREATE TABLE `new_User` (
  `id` text NOT NULL,
  `email` text NOT NULL,
  `phoneNum` text NOT NULL,
  `role` text NOT NULL DEFAULT 'student',
  `referredById` text NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `User_referredById_fkey` FOREIGN KEY (`referredById`) REFERENCES `User` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `User_role_check` CHECK ("role" IN ('ADMIN', 'student', 'TEACHER'))
);
INSERT INTO `new_User` (`id`, `email`, `phoneNum`, `role`, `referredById`) SELECT `id`, `email`, `phoneNum`, `role`, `referredById` FROM `User`;
DROP TABLE `User`;
ALTER TABLE `new_User` RENAME TO `User`;
CREATE UNIQUE INDEX `User_email_key` ON `User` (`email`);
CREATE UNIQUE INDEX `User_email_phoneNum_key` ON `User` (`email`, `phoneNum`);
PRAGMA foreign_keys = on;
