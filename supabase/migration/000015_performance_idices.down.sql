BEGIN;

DROP INDEX IF EXISTS idx_booking_status;
DROP INDEX IF EXISTS idx_booking_booking_time;
DROP INDEX IF EXISTS idx_show_slot_id;
DROP INDEX IF EXISTS idx_show_movie_id;
DROP INDEX IF EXISTS idx_admin_booked_customer_booking_id;

COMMIT;