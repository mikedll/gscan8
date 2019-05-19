DROP TABLE IF EXISTS "gists" CASCADE;
CREATE TABLE "gists" (
	"id" bigserial primary key,
	"user_id" integer,
	"vendor_id" character varying NOT NULL,
  "updated_at" timestamp,
	"title" character varying DEFAULT '' NOT NULL,
	"url" character varying DEFAULT '' NOT NULL,
	"body" character varying DEFAULT '' NOT NULL
);

DROP TABLE IF EXISTS "users" CASCADE;
CREATE TABLE "users" (
	"id" bigserial primary key,
	"credential" character varying DEFAULT '' NOT NULL,
	"authentication_service" character varying DEFAULT '' NOT NULL
);
