BEGIN;

CREATE INDEX IF NOT EXISTS idx_booking_status ON booking (status);
CREATE INDEX IF NOT EXISTS idx_booking_booking_time ON booking (booking_time);
CREATE INDEX IF NOT EXISTS idx_show_slot_id ON show (slot_id);
CREATE INDEX IF NOT EXISTS idx_show_movie_id ON show (movie_id);
CREATE INDEX IF NOT EXISTS idx_admin_booked_customer_booking_id ON admin_booked_customer (booking_id);

COMMENT ON INDEX idx_booking_status IS 'Improves performance when filtering bookings by status';
COMMENT ON INDEX idx_booking_booking_time IS 'Improves performance for timeframe-based queries';
COMMENT ON INDEX idx_show_slot_id IS 'Improves performance when filtering shows by slot';
COMMENT ON INDEX idx_show_movie_id IS 'Improves performance when filtering shows by movie';
COMMENT ON INDEX idx_admin_booked_customer_booking_id IS 'Covers foreign key fk_admin_customer_booking for better join performance';

COMMIT;