-- 一時カラムに変換後のUUIDを入れる
ALTER TABLE messages
    ADD COLUMN new_room_id UUID;

-- 既存データの文字列 → UUIDに変換（仮に既にUUIDが入ってる前提ならそのまま）
UPDATE messages SET new_room_id = room_id::UUID;

-- 古いカラム削除・リネーム
ALTER TABLE messages DROP COLUMN room_id;
ALTER TABLE messages RENAME COLUMN new_room_id TO room_id;

-- 外部キー制約を付ける
ALTER TABLE messages
    ADD CONSTRAINT fk_messages_rooms
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE;
