CREATE TABLE "users" (
  "id" varchar NOT NULL PRIMARY KEY,
  "username" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "email" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "users" ADD CONSTRAINT "username_key" UNIQUE ("username");
ALTER TABLE "users" ADD CONSTRAINT "email_key" UNIQUE ("email");
CREATE INDEX ON "users" ("username");
CREATE INDEX ON "users" ("email");


CREATE TABLE "user_relate_url" (
    "id" bigserial NOT NULL PRIMARY KEY,
    "user_id" varchar NOT NULL,
    "short_url" varchar NOT NULL,
    "origin_url" varchar NOT NULL,
    "status" int NOT NULL DEFAULT 0,
    "expire_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "user_relate_url" ADD CONSTRAINT "short_url_key" UNIQUE ("short_url");
CREATE INDEX ON "user_relate_url" ("user_id");
CREATE INDEX ON "user_relate_url" ("short_url");