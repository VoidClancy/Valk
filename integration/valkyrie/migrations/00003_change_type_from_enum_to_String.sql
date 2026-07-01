-- +goose Up
PRAGMA foreign_keys = off;
CREATE TABLE `new_Comment` (
  `id` text NOT NULL,
  `textify` text NOT NULL,
  `dummy3` integer NOT NULL,
  `dummy1` text NOT NULL,
  `dummy2` text NOT NULL,
  `postId` text NOT NULL,
  `authorId` text NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `Comment_postId_fkey` FOREIGN KEY (`postId`) REFERENCES `Post` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `Comment_authorId_fkey` FOREIGN KEY (`authorId`) REFERENCES `User` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
INSERT INTO `new_Comment` (`id`, `textify`, `dummy3`, `dummy1`, `dummy2`, `postId`, `authorId`) SELECT `id`, `textify`, `dummy3`, `dummy1`, `dummy2`, `postId`, `authorId` FROM `Comment`;
DROP TABLE `Comment`;
ALTER TABLE `new_Comment` RENAME TO `Comment`;
PRAGMA foreign_keys = on;

-- +goose Down
PRAGMA foreign_keys = off;
CREATE TABLE `new_Comment` (
  `id` text NOT NULL,
  `textify` text NOT NULL,
  `dummy3` integer NOT NULL,
  `dummy1` text NOT NULL,
  `dummy2` text NOT NULL,
  `postId` text NOT NULL,
  `authorId` text NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `Comment_authorId_fkey` FOREIGN KEY (`authorId`) REFERENCES `User` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `Comment_postId_fkey` FOREIGN KEY (`postId`) REFERENCES `Post` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `Comment_textify_check` CHECK ("textify" IN ('ADMIN', 'student', 'TEACHER'))
);
INSERT INTO `new_Comment` (`id`, `textify`, `dummy3`, `dummy1`, `dummy2`, `postId`, `authorId`) SELECT `id`, `textify`, `dummy3`, `dummy1`, `dummy2`, `postId`, `authorId` FROM `Comment`;
DROP TABLE `Comment`;
ALTER TABLE `new_Comment` RENAME TO `Comment`;
PRAGMA foreign_keys = on;
