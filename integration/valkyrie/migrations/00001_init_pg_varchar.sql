-- +goose Up
CREATE TYPE "public"."user_roles" AS ENUM ('ADMIN', 'student', 'TEACHER');
CREATE TABLE "public"."User" (
  "id" text NOT NULL,
  "email" text NOT NULL,
  "phoneNum" text NOT NULL,
  "role" "public"."user_roles" NOT NULL DEFAULT 'student',
  "referredById" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "User_referredById_fkey" FOREIGN KEY ("referredById") REFERENCES "public"."User" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE "public"."Post" (
  "id" text NOT NULL,
  "title" VARCHAR(22) NOT NULL,
  "content" text NULL,
  "published" boolean NOT NULL DEFAULT FALSE,
  "authorId" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "Post_authorId_fkey" FOREIGN KEY ("authorId") REFERENCES "public"."User" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE "public"."Category" (
  "id" serial NOT NULL,
  "name" text NOT NULL,
  PRIMARY KEY ("id")
);
CREATE TABLE "public"."CategoryToPost" (
  "postId" text NOT NULL,
  "categoryId" integer NOT NULL,
  PRIMARY KEY ("postId", "categoryId"),
  CONSTRAINT "CategoryToPost_postId_fkey" FOREIGN KEY ("postId") REFERENCES "public"."Post" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "CategoryToPost_categoryId_fkey" FOREIGN KEY ("categoryId") REFERENCES "public"."Category" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE "public"."Comment" (
  "id" text NOT NULL,
  "textify" integer NOT NULL,
  "dummy3" text NOT NULL,
  "dummy1" integer NOT NULL,
  "dummy2" text NOT NULL,
  "postId" text NOT NULL,
  "authorId" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "Comment_postId_fkey" FOREIGN KEY ("postId") REFERENCES "public"."Post" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "Comment_authorId_fkey" FOREIGN KEY ("authorId") REFERENCES "public"."User" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE "public"."Profile" (
  "id" text NOT NULL,
  "bio" text NULL,
  "userId" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "Profile_userId_fkey" FOREIGN KEY ("userId") REFERENCES "public"."User" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE UNIQUE INDEX "User_email_key" ON "public"."User" ("email");
CREATE UNIQUE INDEX "User_phoneNum_key" ON "public"."User" ("phoneNum");
CREATE UNIQUE INDEX "User_email_phoneNum_key" ON "public"."User" ("email", "phoneNum");
CREATE UNIQUE INDEX "Category_name_key" ON "public"."Category" ("name");
CREATE UNIQUE INDEX "Profile_userId_key" ON "public"."Profile" ("userId");

-- +goose Down
ALTER TABLE "public"."Profile" DROP CONSTRAINT "Profile_userId_fkey";
ALTER TABLE "public"."Post" DROP CONSTRAINT "Post_authorId_fkey";
ALTER TABLE "public"."Comment" DROP CONSTRAINT "Comment_postId_fkey", DROP CONSTRAINT "Comment_authorId_fkey";
ALTER TABLE "public"."CategoryToPost" DROP CONSTRAINT "CategoryToPost_postId_fkey", DROP CONSTRAINT "CategoryToPost_categoryId_fkey";
DROP TABLE "public"."User";
DROP TYPE "public"."user_roles";
DROP TABLE "public"."Profile";
DROP TABLE "public"."Post";
DROP TABLE "public"."Comment";
DROP TABLE "public"."Category";
DROP TABLE "public"."CategoryToPost";
