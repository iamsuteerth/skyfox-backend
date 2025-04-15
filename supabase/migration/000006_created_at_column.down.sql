BEGIN;

ALTER TABLE usertable
DROP COLUMN created_at;

COMMIT;