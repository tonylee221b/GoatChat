-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  username varchar(32) NOT NULL,
  email varchar(254) NOT NULL,
  phone_number varchar(32),
  status varchar(16) NOT NULL DEFAULT 'active',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  CONSTRAINT users_username_not_blank CHECK (length(btrim(username)) > 0),
  CONSTRAINT users_email_not_blank CHECK (length(btrim(email)) > 0),
  CONSTRAINT users_email_shape CHECK (position('@' in email) > 1),
  CONSTRAINT users_status_check CHECK (status IN ('active', 'inactive', 'deleted'))
);

CREATE UNIQUE INDEX users_username_unique_active
  ON users (lower(username))
  WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX users_email_unique_active
  ON users (lower(email))
  WHERE deleted_at IS NULL;

CREATE TABLE auth_credentials (
  user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  password_hash text NOT NULL,
  password_algorithm varchar(32) NOT NULL,
  password_changed_at timestamptz NOT NULL DEFAULT now(),
  failed_login_count integer NOT NULL DEFAULT 0,
  locked_until timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT auth_credentials_password_hash_not_blank CHECK (length(btrim(password_hash)) > 0),
  CONSTRAINT auth_credentials_password_algorithm_check CHECK (password_algorithm IN ('bcrypt', 'argon2id')),
  CONSTRAINT auth_credentials_failed_login_count_check CHECK (failed_login_count >= 0)
);

CREATE TABLE chat_rooms (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  room_type varchar(16) NOT NULL,
  name varchar(100),
  description varchar(500),
  owner_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  last_message_id uuid,
  last_message_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  CONSTRAINT chat_rooms_room_type_check CHECK (room_type IN ('direct', 'group')),
  CONSTRAINT chat_rooms_group_name_required CHECK (
    room_type <> 'group' OR (name IS NOT NULL AND length(btrim(name)) > 0)
  )
);

CREATE INDEX chat_rooms_owner_id_idx ON chat_rooms (owner_id);
CREATE INDEX chat_rooms_last_message_at_idx ON chat_rooms (last_message_at DESC NULLS LAST);

CREATE TABLE chat_room_members (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id uuid NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  role varchar(16) NOT NULL,
  joined_at timestamptz NOT NULL DEFAULT now(),
  left_at timestamptz,
  last_read_message_id uuid,
  last_read_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT chat_room_members_role_check CHECK (role IN ('owner', 'member')),
  CONSTRAINT chat_room_members_left_after_joined CHECK (left_at IS NULL OR left_at >= joined_at)
);

CREATE UNIQUE INDEX chat_room_members_active_member_unique
  ON chat_room_members (room_id, user_id)
  WHERE left_at IS NULL;

CREATE INDEX chat_room_members_user_id_idx ON chat_room_members (user_id);
CREATE INDEX chat_room_members_room_id_idx ON chat_room_members (room_id);

CREATE TABLE messages (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id uuid NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
  sender_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  message_type varchar(16) NOT NULL DEFAULT 'text',
  content text NOT NULL,
  client_message_id varchar(128) NOT NULL,
  status varchar(16) NOT NULL DEFAULT 'sent',
  sent_at timestamptz NOT NULL DEFAULT now(),
  edited_at timestamptz,
  deleted_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT messages_message_type_check CHECK (message_type IN ('text', 'system')),
  CONSTRAINT messages_status_check CHECK (status IN ('sent', 'edited', 'deleted')),
  CONSTRAINT messages_content_not_blank CHECK (length(btrim(content)) > 0),
  CONSTRAINT messages_content_length_check CHECK (char_length(content) <= 4000),
  CONSTRAINT messages_client_message_id_not_blank CHECK (length(btrim(client_message_id)) > 0)
);

CREATE UNIQUE INDEX messages_sender_client_message_unique
  ON messages (sender_id, client_message_id);

CREATE INDEX messages_room_cursor_idx ON messages (room_id, sent_at DESC, id DESC);
CREATE INDEX messages_sender_id_idx ON messages (sender_id);

ALTER TABLE chat_rooms
  ADD CONSTRAINT chat_rooms_last_message_id_fkey
  FOREIGN KEY (last_message_id) REFERENCES messages(id) ON DELETE SET NULL;

ALTER TABLE chat_room_members
  ADD CONSTRAINT chat_room_members_last_read_message_id_fkey
  FOREIGN KEY (last_read_message_id) REFERENCES messages(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE chat_room_members
  DROP CONSTRAINT IF EXISTS chat_room_members_last_read_message_id_fkey;

ALTER TABLE chat_rooms
  DROP CONSTRAINT IF EXISTS chat_rooms_last_message_id_fkey;

DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chat_room_members;
DROP TABLE IF EXISTS chat_rooms;
DROP TABLE IF EXISTS auth_credentials;
DROP TABLE IF EXISTS users;
