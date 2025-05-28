BEGIN;

DROP INDEX IF EXISTS idx_booking_seat_mapping_seat_number;
DROP INDEX IF EXISTS idx_stafftable_username; 

CREATE INDEX idx_booking_date_status ON booking(date, status) 
WHERE status IN ('Pending', 'Confirmed');
COMMENT ON INDEX idx_booking_date_status IS 'Optimizes booking queries filtered by date and active status';

CREATE INDEX idx_show_date_slot ON show(date, slot_id);
COMMENT ON INDEX idx_show_date_slot IS 'Improves show availability queries by date and time slot';

CREATE INDEX idx_wallet_transaction_user_time ON wallet_transaction(username, timestamp DESC);
COMMENT ON INDEX idx_wallet_transaction_user_time IS 'Optimizes wallet transaction history queries';

CREATE INDEX idx_payment_transaction_booking_status ON payment_transaction(booking_id, status);
COMMENT ON INDEX idx_payment_transaction_booking_status IS 'Improves payment status lookup performance';

CREATE INDEX idx_active_bookings ON booking(show_id, customer_username) 
WHERE status IN ('Pending', 'Confirmed');
COMMENT ON INDEX idx_active_bookings IS 'Optimizes queries for active bookings only';

CREATE INDEX idx_active_reset_tokens ON password_reset_tokens(email, expires_at) 
WHERE used = FALSE;
COMMENT ON INDEX idx_active_reset_tokens IS 'Improves password reset token lookup for unused tokens';

CREATE INDEX idx_wallet_balance_lookup ON customer_wallet(username, balance) 
WHERE balance > 0;
COMMENT ON INDEX idx_wallet_balance_lookup IS 'Optimizes wallet balance queries for customers with funds';

CREATE INDEX idx_wallet_txn_aggregation ON wallet_transaction(username, transaction_type, amount);
COMMENT ON INDEX idx_wallet_txn_aggregation IS 'Improves wallet transaction aggregation queries';

CREATE INDEX idx_booking_details_covering ON booking(id, show_id, customer_username, status, amount_paid, booking_time);
COMMENT ON INDEX idx_booking_details_covering IS 'Covering index for booking detail queries to avoid table lookups';

CREATE INDEX idx_show_movie_date ON show(movie_id, date, slot_id);
COMMENT ON INDEX idx_show_movie_date IS 'Optimizes movie show schedule queries';

COMMIT;
