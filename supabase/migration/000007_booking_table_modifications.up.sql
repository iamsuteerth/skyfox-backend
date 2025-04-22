BEGIN;

ALTER TABLE booking
ADD COLUMN customer_username VARCHAR(30) NULL,
ADD CONSTRAINT fk_booking_customer_username 
    FOREIGN KEY (customer_username) REFERENCES customertable(username) 
    ON DELETE CASCADE;

ALTER TABLE booking 
ADD CONSTRAINT check_customer_type 
    CHECK (
        (customer_id IS NOT NULL AND customer_username IS NULL) OR 
        (customer_id IS NULL AND customer_username IS NOT NULL) 
    );

COMMIT;