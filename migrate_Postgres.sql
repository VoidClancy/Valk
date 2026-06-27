-- +goose Up
CREATE TYPE "user_roles" AS ENUM (
  'ADMIN',
  'student',
  'TEACHER'
);

CREATE TYPE "Clancy" AS ENUM (
  'VOICES',
  'BLEH'
);

CREATE TABLE "User" (
  "twistedID" TEXT NOT NULL,
  "email" TEXT NOT NULL,
  "phoneNum" TEXT NOT NULL,
  "role" "user_roles" NOT NULL DEFAULT 'student',
  "referredById" TEXT NULL,
  CONSTRAINT "User_pkey" PRIMARY KEY ("twistedID"),
  CONSTRAINT "User_email_key" UNIQUE ("email"),
  CONSTRAINT "User_email_phoneNum_key" UNIQUE ("email", "phoneNum"),
  CONSTRAINT "User_referredById_fkey" FOREIGN KEY ("referredById") REFERENCES "User" ("twistedID")
);

CREATE TABLE "Profile" (
  "id" TEXT NOT NULL,
  "bio" TEXT NULL,
  "userId" TEXT NOT NULL,
  CONSTRAINT "Profile_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "Profile_userId_key" UNIQUE ("userId"),
  CONSTRAINT "Profile_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User" ("twistedID")
);

CREATE TABLE "Post" (
  "id" TEXT NOT NULL,
  "title" TEXT NOT NULL,
  "content" TEXT NULL,
  "published" BOOLEAN NOT NULL DEFAULT FALSE,
  "authorId" TEXT NOT NULL,
  CONSTRAINT "Post_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "Post_authorId_fkey" FOREIGN KEY ("authorId") REFERENCES "User" ("twistedID")
);

CREATE TABLE "Comment" (
  "id" TEXT NOT NULL,
  "text" TEXT NOT NULL,
  "postId" TEXT NOT NULL,
  "authorId" TEXT NOT NULL,
  CONSTRAINT "Comment_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "Comment_postId_fkey" FOREIGN KEY ("postId") REFERENCES "Post" ("id"),
  CONSTRAINT "Comment_authorId_fkey" FOREIGN KEY ("authorId") REFERENCES "User" ("twistedID")
);

CREATE TABLE "Category" (
  "id" SERIAL NOT NULL,
  "name" TEXT NOT NULL,
  CONSTRAINT "Category_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "Category_name_key" UNIQUE ("name")
);

CREATE TABLE "CategoryToPost" (
  "postId" TEXT NOT NULL,
  "categoryId" INTEGER NOT NULL,
  CONSTRAINT "CategoryToPost_pkey" PRIMARY KEY ("postId", "categoryId"),
  CONSTRAINT "CategoryToPost_postId_fkey" FOREIGN KEY ("postId") REFERENCES "Post" ("id"),
  CONSTRAINT "CategoryToPost_categoryId_fkey" FOREIGN KEY ("categoryId") REFERENCES "Category" ("id")
);


-- +goose Down
DROP TABLE IF EXISTS "CategoryToPost";
DROP TABLE IF EXISTS "Category";
DROP TABLE IF EXISTS "Comment";
DROP TABLE IF EXISTS "Post";
DROP TABLE IF EXISTS "Profile";
DROP TABLE IF EXISTS "User";
DROP TYPE IF EXISTS "Clancy";
DROP TYPE IF EXISTS "user_roles";
