BEGIN;

DROP INDEX IF EXISTS idx_password_reset_tokens_token;
DROP INDEX IF EXISTS idx_password_reset_tokens_email;
DROP TABLE IF EXISTS password_reset_tokens;

COMMIT;