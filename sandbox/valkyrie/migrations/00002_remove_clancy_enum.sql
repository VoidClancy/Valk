-- +goose Up
DROP TYPE "public"."Clancy";

-- +goose Down
CREATE TYPE "public"."Clancy" AS ENUM ('VOICES', 'BLEH');
