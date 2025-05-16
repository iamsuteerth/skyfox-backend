BEGIN;

CREATE TYPE wallet_transaction_type AS ENUM ('ADD', 'DEDUCT');

COMMIT;