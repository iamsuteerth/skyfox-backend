BEGIN;

ALTER TABLE booking DROP CONSTRAINT IF EXISTS unique_customer_id; 

ALTER TABLE admin_booked_customer DROP CONSTRAINT IF EXISTS fk_admin_customer_booking;

ALTER TABLE booking
ADD CONSTRAINT fk_booking_customer
FOREIGN KEY (customer_id) 
REFERENCES admin_booked_customer(id)
ON DELETE NO ACTION;

COMMIT;