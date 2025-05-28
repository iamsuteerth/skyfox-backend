BEGIN;

ALTER TABLE wallet_transaction DROP CONSTRAINT fk_wallet_username;
ALTER TABLE customer_wallet DROP CONSTRAINT fk_customer_wallet_username;

ALTER TABLE wallet_transaction 
ADD CONSTRAINT fk_wallet_username 
FOREIGN KEY (username) REFERENCES customertable(username) ON DELETE CASCADE;

ALTER TABLE customer_wallet 
ADD CONSTRAINT fk_customer_wallet_username 
FOREIGN KEY (username) REFERENCES customertable(username) ON DELETE CASCADE;

COMMIT;