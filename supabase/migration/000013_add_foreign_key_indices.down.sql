BEGIN;

DROP INDEX IF EXISTS idx_booking_customer_username;
DROP INDEX IF EXISTS idx_booking_show_id;
DROP INDEX IF EXISTS idx_booking_seat_mapping_booking_id;
DROP INDEX IF EXISTS idx_booking_seat_mapping_seat_number;
DROP INDEX IF EXISTS idx_stafftable_username;

COMMIT;