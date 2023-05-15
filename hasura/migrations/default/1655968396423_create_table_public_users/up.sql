CREATE TABLE "public"."users" ("id" serial NOT NULL, "name" text, "created_at" timestamptz, "updated_at" timestamptz, PRIMARY KEY ("id") , UNIQUE ("id"));
alter table "public"."users" alter column "created_at" set default now();
alter table "public"."users" alter column "updated_at" set default now();
