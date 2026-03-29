table "files" {
  schema = schema.public
  column "id" {
    null = false
    type = text
  }
  column "file_name" {
    null = false
    type = text
  }
  column "storage_key" {
    null = true
    type = text
  }
  column "mime_type" {
    null = false
    type = text
  }
  column "file_size" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "format" {
    null = true
    type = text
  }
  column "width" {
    null = true
    type = bigint
  }
  column "height" {
    null = true
    type = bigint
  }
  column "user_id" {
    null = true
    type = text
  }
  column "thumbnail_storage_key" {
    null = true
    type = text
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_files_thumbnail_storage_key" {
    unique  = true
    columns = [column.thumbnail_storage_key]
  }
  index "idx_files_user_id" {
    columns = [column.user_id]
  }
}
table "sessions" {
  schema = schema.public
  column "id" {
    null = false
    type = text
  }
  column "token_hash" {
    null = false
    type = text
  }
  column "access_jti" {
    null = false
    type = text
  }
  column "user_id" {
    null = false
    type = text
  }
  column "expires_at" {
    null = false
    type = timestamptz
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_sessions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = CASCADE
    on_delete   = SET_NULL
  }
}
table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = text
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = true
    type = text
  }
  column "password" {
    null = false
    type = text
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "deleted_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_deleted_at" {
    columns = [column.deleted_at]
  }
  index "idx_users_email" {
    unique  = true
    columns = [column.email]
  }
}
schema "public" {
  comment = "standard public schema"
}
