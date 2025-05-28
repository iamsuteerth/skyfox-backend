BEGIN;

DROP INDEX IF EXISTS idx_booking_date_status;
DROP INDEX IF EXISTS idx_show_date_slot;
DROP INDEX IF EXISTS idx_wallet_transaction_user_time;
DROP INDEX IF EXISTS idx_payment_transaction_booking_status;
DROP INDEX IF EXISTS idx_active_bookings;
DROP INDEX IF EXISTS idx_active_reset_tokens;
DROP INDEX IF EXISTS idx_wallet_balance_lookup;
DROP INDEX IF EXISTS idx_wallet_txn_aggregation;
DROP INDEX IF EXISTS idx_booking_details_covering;
DROP INDEX IF EXISTS idx_show_movie_date;

CREATE INDEX idx_booking_seat_mapping_seat_number ON booking_seat_mapping(seat_number);
CREATE INDEX idx_stafftable_username ON stafftable(username);

COMMIT;
