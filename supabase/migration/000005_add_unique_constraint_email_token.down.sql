BEGIN;

ALTER TABLE password_reset_tokens 
DROP CONSTRAINT unique_email_token;

COMMIT;
