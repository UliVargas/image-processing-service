-- Create "files" table
CREATE TABLE "files" (
  "id" text NOT NULL,
  "file_name" text NOT NULL,
  "storage_key" text NULL,
  "thumbnail_storage_key" text NULL,
  "mime_type" text NOT NULL,
  "file_size" bigint NOT NULL,
  "user_id" text NULL,
  "format" text NOT NULL,
  "width" bigint NOT NULL,
  "height" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_files_thumbnail_storage_key" to table: "files"
CREATE UNIQUE INDEX "idx_files_thumbnail_storage_key" ON "files" ("thumbnail_storage_key");
-- Create index "idx_files_user_id" to table: "files"
CREATE INDEX "idx_files_user_id" ON "files" ("user_id");
-- Create "users" table
CREATE TABLE "users" (
  "id" text NOT NULL,
  "name" text NOT NULL,
  "email" text NULL,
  "password" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Create "sessions" table
CREATE TABLE "sessions" (
  "id" text NOT NULL,
  "token_hash" text NOT NULL,
  "access_jti" text NOT NULL,
  "user_id" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_sessions_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
