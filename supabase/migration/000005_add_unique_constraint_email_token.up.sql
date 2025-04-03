BEGIN;

ALTER TABLE password_reset_tokens 
ADD CONSTRAINT unique_email_token UNIQUE (email, token);

COMMIT;
