CREATE TABLE "public"."statuses" ("id" serial NOT NULL, "name" text, "created_at" timestamptz, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id") , UNIQUE ("id"));
alter table "public"."statuses" alter column "created_at" set default now();
alter table "public"."statuses" alter column "updated_at" set default now();
