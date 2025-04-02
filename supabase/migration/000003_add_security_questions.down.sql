BEGIN;

DROP INDEX IF EXISTS idx_security_question_id;

ALTER TABLE customertable 
DROP COLUMN IF EXISTS security_question_id,
DROP COLUMN IF EXISTS security_answer_hash;

DROP TABLE IF EXISTS security_questions;

COMMIT;