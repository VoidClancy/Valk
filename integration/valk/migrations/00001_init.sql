-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS hstore;
CREATE EXTENSION IF NOT EXISTS ltree;

CREATE TYPE "public"."user_roles" AS ENUM ('ADMIN', 'student', 'TEACHER');
CREATE TABLE "public"."User" (
  "id" text NOT NULL,
  "email" text NOT NULL,
  "phoneNum" text NOT NULL,
  "password" text NULL,
  "role" "public"."user_roles" NOT NULL DEFAULT 'student',
  "roleOptional" "public"."user_roles" NULL,
  "loginCount" integer NOT NULL DEFAULT 0,
  "referredById" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "User_referredById_fkey" FOREIGN KEY ("referredById") REFERENCES "public"."User" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE "public"."DefaultsTest" (
  "uuid4" text NOT NULL,
  "uuid7" text NOT NULL,
  "uuidNoArgs" text NOT NULL,
  "cuid1" text NOT NULL,
  "cuid2" text NOT NULL,
  "cuidNoArgs" text NOT NULL,
  "ulid" text NOT NULL,
  "nanoid" text NOT NULL,
  "now" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("uuid4")
);
CREATE TABLE "public"."AllFieldsSoFar" (
  "id" serial NOT NULL,
  "stringReq" text NOT NULL,
  "stringOpt" text NULL,
  "stringDefault" text NOT NULL DEFAULT 'default',
  "stringVarchar" character varying(255) NOT NULL,
  "stringChar" character(10) NOT NULL,
  "bitVal" bit(10) NOT NULL,
  "varBitVal" bit varying NOT NULL,
  "inetVal" inet NOT NULL,
  "xmlVal" xml NOT NULL,
  "cuidDefault" text NOT NULL,
  "cuid1Default" text NOT NULL,
  "cuid2Default" text NOT NULL,
  "uuidDefault" text NOT NULL,
  "uuid4Default" text NOT NULL,
  "uuid7Default" text NOT NULL,
  "ulidDefault" text NOT NULL,
  "nanoidDefault" text NOT NULL,
  "uuidDb" uuid NOT NULL,
  "intReq" integer NOT NULL,
  "intOpt" integer NULL,
  "intDefault" integer NOT NULL DEFAULT 42,
  "integerVal" integer NOT NULL,
  "smallInt" smallint NOT NULL,
  "tinyInt" integer NOT NULL,
  "oidVal" oid NOT NULL,
  "bigIntReq" bigint NOT NULL,
  "bigIntOpt" bigint NULL,
  "floatReq" double precision NOT NULL,
  "floatOpt" double precision NULL,
  "realVal" real NOT NULL,
  "decimalReq" numeric NOT NULL,
  "decimalOpt" numeric NULL,
  "decimalPrecise" numeric(10,2) NOT NULL,
  "moneyVal" money NOT NULL,
  "boolReq" boolean NOT NULL,
  "boolOpt" boolean NULL,
  "boolDefault" boolean NOT NULL DEFAULT FALSE,
  "dateTimeReq" timestamp NOT NULL,
  "dateTimeOpt" timestamp NULL,
  "dateTimeDefault" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updatedAt" timestamp NOT NULL,
  "dateTimeTz" timestamptz NOT NULL,
  "timestampVal" timestamp(3) NOT NULL,
  "timeVal" time NOT NULL,
  "timetzVal" timetz NOT NULL,
  "jsonReq" jsonb NOT NULL,
  "jsonOpt" jsonb NULL,
  "jsonVal" json NOT NULL,
  "bytesReq" bytea NOT NULL,
  "bytesOpt" bytea NULL,
  "hstoreField" hstore NULL,
  "ltreeField" ltree NOT NULL,
  "citextField" citext NULL,
  PRIMARY KEY ("id")
);
CREATE TABLE "public"."Post" (
  "id" text NOT NULL,
  "title" text NOT NULL,
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
  "meta" jsonb NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "Comment_postId_fkey" FOREIGN KEY ("postId") REFERENCES "public"."Post" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "Comment_authorId_fkey" FOREIGN KEY ("authorId") REFERENCES "public"."User" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE "public"."Profile" (
  "id" text NOT NULL,
  "bio" text NULL,
  "userId" text NOT NULL,
  "createdAt" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "Profile_userId_fkey" FOREIGN KEY ("userId") REFERENCES "public"."User" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE UNIQUE INDEX "User_email_key" ON "public"."User" ("email");
CREATE UNIQUE INDEX "User_phoneNum_key" ON "public"."User" ("phoneNum");
CREATE UNIQUE INDEX "emailPhone" ON "public"."User" ("email", "phoneNum");
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
DROP TABLE "public"."DefaultsTest";
DROP TABLE "public"."AllFieldsSoFar";
