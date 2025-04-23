BEGIN;

ALTER TABLE admin_booked_customer DROP CONSTRAINT IF EXISTS fk_admin_customer_booking;

ALTER TABLE admin_booked_customer DROP COLUMN IF EXISTS booking_id;

ALTER TABLE booking ADD CONSTRAINT unique_customer_id UNIQUE (customer_id);

ALTER TABLE admin_booked_customer
ADD CONSTRAINT fk_admin_customer_booking
FOREIGN KEY (id)
REFERENCES booking(customer_id)
ON DELETE CASCADE;

COMMIT;