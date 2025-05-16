BEGIN;

CREATE TABLE wallet_transaction (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL,
    username VARCHAR(30) NOT NULL,
    booking_id BIGINT, -- nullable
    transaction_id VARCHAR(255) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    transaction_type wallet_transaction_type NOT NULL,
    CONSTRAINT fk_wallet_id FOREIGN KEY (wallet_id) REFERENCES customer_wallet(id) ON DELETE CASCADE,
    CONSTRAINT fk_wallet_username FOREIGN KEY (username) REFERENCES customertable(username),
    CONSTRAINT fk_wallet_booking FOREIGN KEY (booking_id) REFERENCES booking(id) ON DELETE SET NULL
);

CREATE INDEX idx_wallet_transaction_wallet_id ON wallet_transaction(wallet_id);
CREATE INDEX idx_wallet_transaction_username ON wallet_transaction(username);
CREATE INDEX idx_wallet_transaction_booking_id ON wallet_transaction(booking_id);
CREATE INDEX idx_wallet_transaction_timestamp ON wallet_transaction(timestamp);

COMMIT;