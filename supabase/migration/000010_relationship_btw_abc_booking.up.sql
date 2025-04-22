BEGIN;

ALTER TABLE booking ADD CONSTRAINT unique_customer_id UNIQUE (customer_id);

ALTER TABLE booking DROP CONSTRAINT IF EXISTS fk_booking_customer;

ALTER TABLE admin_booked_customer
ADD CONSTRAINT fk_admin_customer_booking
FOREIGN KEY (id)
REFERENCES booking(customer_id)
ON DELETE CASCADE;

COMMIT;