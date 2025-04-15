-- 拡張機能（UUID生成に必要）
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ユーザーテーブル
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- チャットルームテーブル
CREATE TABLE rooms (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT,
  display_name TEXT,
  is_group BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_message TEXT
);

-- ルームメンバーテーブル
CREATE TABLE room_members (
  room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (room_id, user_id)
);

-- メッセージテーブル（room_idは最初からUUIDとして作成）
CREATE TABLE messages (
  id SERIAL PRIMARY KEY,
  room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  sender_id INTEGER NOT NULL,
  sender TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ルームの既読情報テーブル
CREATE TABLE room_reads (
  user_id INT NOT NULL,
  room_id UUID NOT NULL,
  last_read_at TIMESTAMP NOT NULL,
  PRIMARY KEY (user_id, room_id)
);
