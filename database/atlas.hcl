table "account" {
  schema = schema.public
  column "uuid" {
    null = false
    type = uuid
  }
  column "id" {
    null    = true
    type    = text
    comment = "The unique identifier for the account within the system."
  }
  column "name" {
    null    = true
    type    = text
    comment = "The human readable / display name of the account"
  }
  column "owner_uuid" {
    null = true
    type = uuid
  }
  column "flags" {
    null = true
    type = jsonb
  }
  column "created" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.uuid]
  }
  foreign_key "owner_uuid" {
    columns     = [column.owner_uuid]
    ref_columns = [table.user.column.uuid]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  index "unique_account_slug" {
    unique  = true
    columns = [column.id]
  }
}
table "api_key" {
  schema = schema.public
  column "key" {
    null    = false
    type    = text
    comment = "The one-way hashed API key"
  }
  column "name" {
    null    = false
    type    = text
    comment = "User-provided name"
  }
  column "user_uuid" {
    null = false
    type = uuid
  }
  column "issued" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "expires" {
    null = true
    type = timestamp
  }
  column "last_used" {
    null = true
    type = timestamp
  }
  column "decription" {
    null    = true
    type    = text
    comment = "User provided description"
  }
  foreign_key "api_key_user_uuid" {
    columns     = [column.user_uuid]
    ref_columns = [table.user.column.uuid]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
}
table "login_session" {
  schema = schema.public
  column "session_id" {
    null = false
    type = uuid
  }
  column "key" {
    null = false
    type = text
  }
  column "value" {
    null = true
    type = text
  }
  primary_key {
    columns = [column.session_id, column.key]
  }
}
table "oauth_login" {
  schema = schema.public
  column "provider" {
    null = false
    type = text
  }
  column "provider_id" {
    null = false
    type = text
  }
  column "provider_username" {
    null = true
    type = text
  }
  column "access_token" {
    null = false
    type = text
  }
  column "refresh_token" {
    null = false
    type = text
  }
  column "authorised_user" {
    null = false
    type = uuid
  }
  primary_key {
    columns = [column.provider, column.provider_id]
  }
  foreign_key "authorised_user_uuid" {
    columns     = [column.authorised_user]
    ref_columns = [table.user.column.uuid]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}
table "pod_services" {
  schema = schema.public
  column "uuid" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "pod_uuid" {
    null = true
    type = uuid
  }
  column "name" {
    null = true
    type = text
  }
  column "port" {
    null = false
    type = smallint
  }
  column "target_port" {
    null = true
    type = smallint
  }
  column "protocol" {
    null = true
    type = character_varying(4)
  }
  column "domain_name" {
    null = true
    type = text
  }
  primary_key {
    columns = [column.uuid]
  }
  foreign_key "pod_uuid_fk" {
    columns     = [column.pod_uuid]
    ref_columns = [table.project_pods.column.uuid]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  index "unique_domain_name" {
    unique  = true
    columns = [column.domain_name]
  }
}
table "pod_volumes" {
  schema = schema.public
  column "uuid" {
    null = false
    type = uuid
  }
  column "pod_uuid" {
    null = false
    type = uuid
  }
  column "name" {
    null = false
    type = text
  }
  column "size" {
    null = false
    type = smallint
  }
  column "mount_path" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.uuid]
  }
  foreign_key "pod_uuid_fk" {
    columns     = [column.pod_uuid]
    ref_columns = [table.project_pods.column.uuid]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
}
table "policy" {
  schema = schema.public
  column "uuid" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "account_uuid" {
    null    = true
    type    = uuid
    comment = "The account to which this policy belongs"
  }
  column "id" {
    null = true
    type = text
  }
  column "current_version" {
    null = true
    type = smallint
  }
  column "content" {
    null    = true
    type    = text
    comment = "The content of the active revision of the policy"
  }
  column "date_added" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "added_by" {
    null    = true
    type    = uuid
    comment = "UUID of the user who added this"
  }
  column "notes" {
    null    = true
    type    = text
    comment = "Space for user to write some notes about this policy"
  }
  primary_key {
    columns = [column.uuid]
  }
  foreign_key "account_uuid_fk" {
    columns     = [column.account_uuid]
    ref_columns = [table.account.column.uuid]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
}
table "policy_attachment" {
  schema = schema.public
  column "account_uuid" {
    null = false
    type = uuid
  }
  column "requestor_id" {
    null    = false
    type    = text
    comment = "The resource (actor) that is requesting the action."
  }
  column "policy_uuid" {
    null = false
    type = uuid
  }
  column "date_attached" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "valid_until" {
    null = true
    type = timestamp
  }
  column "attached_by" {
    null    = true
    type    = text
    comment = "The user id that originally attached this policy"
  }
  primary_key {
    columns = [column.account_uuid, column.requestor_id, column.policy_uuid]
  }
  foreign_key "account_uuid" {
    columns     = [column.account_uuid]
    ref_columns = [table.account.column.uuid]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  foreign_key "policy_uuid" {
    columns     = [column.policy_uuid]
    ref_columns = [table.policy.column.uuid]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
}
table "policy_version" {
  schema = schema.public
  column "uuid" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "policy_uuid" {
    null = true
    type = uuid
  }
  column "version_number" {
    null = false
    type = smallint
  }
  column "content" {
    null    = true
    type    = text
    comment = "The policy document itself"
  }
  column "comment" {
    null    = false
    type    = text
    comment = "Commit message for the revision"
  }
  column "user_uuid" {
    null    = true
    type    = uuid
    comment = "User who made the revision"
  }
  column "date_made" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.uuid]
  }
  foreign_key "policy_uuid" {
    columns     = [column.policy_uuid]
    ref_columns = [table.policy.column.uuid]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
}
table "project" {
  schema = schema.public
  column "uuid" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "id" {
    null    = false
    type    = text
    comment = "Unique identifier for the project within the user's account"
  }
  column "name" {
    null    = true
    type    = text
    comment = "Human readable / display name for the project"
  }
  column "account_uuid" {
    null = true
    type = uuid
  }
  column "created" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.uuid]
  }
  foreign_key "account_fk" {
    columns     = [column.account_uuid]
    ref_columns = [table.account.column.uuid]
    on_update   = CASCADE
    on_delete   = SET_NULL
  }
  index "unique_project_slug_per_account" {
    unique  = true
    columns = [column.account_uuid, column.id]
  }
}
table "project_pods" {
  schema = schema.public
  column "uuid" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "id" {
    null    = true
    type    = text
    comment = "The unique name for the deployment in kubernetes, used as the kuberenetes name."
  }
  column "name" {
    null    = true
    type    = text
    comment = "Human readable / display name for the pod"
  }
  column "image" {
    null    = true
    type    = text
    comment = "The OCI image for the pod to run"
  }
  column "tag" {
    null    = true
    type    = text
    comment = "The image tag to run"
  }
  column "project_uuid" {
    null = true
    type = uuid
  }
  column "environment" {
    null = true
    type = jsonb
  }
  primary_key {
    columns = [column.uuid]
  }
  foreign_key "project_fk" {
    columns     = [column.project_uuid]
    ref_columns = [table.project.column.uuid]
    on_update   = CASCADE
    on_delete   = SET_NULL
  }
}
table "settings" {
  schema = schema.public
  column "section" {
    null = false
    type = text
  }
  column "key" {
    null = false
    type = text
  }
  column "value" {
    null = true
    type = text
  }
  primary_key {
    columns = [column.section, column.key]
  }

}
table "user" {
  schema = schema.public
  column "uuid" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "main_provider" {
    null    = true
    type    = text
    comment = "The provider string for this username, eg github, gitlab, podinate"
  }
  column "id" {
    null = false
    type = text
  }
  column "display_name" {
    null    = true
    type    = text
    comment = "The user's human name"
  }
  column "password_hash" {
    null = true
    type = text
  }
  column "avatar_url" {
    null = true
    type = text
  }
  column "created" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "email" {
    null = false
    type = text
  }
  column "flags" {
    null = true
    type = jsonb
  }
  primary_key {
    columns = [column.uuid]
  }
}
schema "public" {
  comment = "standard public schema"
}
