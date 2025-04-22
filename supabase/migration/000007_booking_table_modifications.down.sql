BEGIN;

ALTER TABLE booking DROP CONSTRAINT check_customer_type;
ALTER TABLE booking DROP CONSTRAINT fk_booking_customer_username;
ALTER TABLE booking DROP COLUMN customer_username;

COMMIT;
