-- +goose Up
DROP INDEX "public"."User_email_phoneNums_key";
ALTER TABLE "public"."User" RENAME COLUMN "phoneNums" TO "phoneNum";
CREATE UNIQUE INDEX "User_email_phoneNum_key" ON "public"."User" ("email", "phoneNum");

-- +goose Down
DROP INDEX "public"."User_email_phoneNum_key";
ALTER TABLE "public"."User" RENAME COLUMN "phoneNum" TO "phoneNums";
CREATE UNIQUE INDEX "User_email_phoneNums_key" ON "public"."User" ("email", "phoneNums");
