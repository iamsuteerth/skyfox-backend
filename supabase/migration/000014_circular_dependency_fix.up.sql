BEGIN;

ALTER TABLE admin_booked_customer DROP CONSTRAINT IF EXISTS fk_admin_customer_booking;
ALTER TABLE booking DROP CONSTRAINT IF EXISTS unique_customer_id;

ALTER TABLE admin_booked_customer ADD COLUMN booking_id BIGINT;

ALTER TABLE admin_booked_customer 
ADD CONSTRAINT fk_admin_customer_booking
FOREIGN KEY (booking_id) REFERENCES booking(id) ON DELETE CASCADE;

COMMIT;