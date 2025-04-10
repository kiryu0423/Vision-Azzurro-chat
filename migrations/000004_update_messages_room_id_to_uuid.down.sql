ALTER TABLE messages DROP CONSTRAINT fk_messages_rooms;

ALTER TABLE messages
    ADD COLUMN old_room_id TEXT;

UPDATE messages SET old_room_id = room_id::TEXT;

ALTER TABLE messages DROP COLUMN room_id;
ALTER TABLE messages RENAME COLUMN old_room_id TO room_id;
