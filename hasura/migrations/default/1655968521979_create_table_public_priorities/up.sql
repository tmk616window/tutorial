CREATE TABLE "public"."priorities" ("id" serial NOT NULL, "name" text, "created_at" timestamptz, "updated_at" timestamptz, PRIMARY KEY ("id") );
alter table "public"."priorities" alter column "created_at" set default now();
alter table "public"."priorities" alter column "updated_at" set default now();
