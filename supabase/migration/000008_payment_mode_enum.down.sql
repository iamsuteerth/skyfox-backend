BEGIN;

ALTER TABLE booking ALTER COLUMN payment_type DROP DEFAULT;

ALTER TABLE booking 
ALTER COLUMN payment_type TYPE VARCHAR(50) 
USING payment_type::VARCHAR(50);

ALTER TABLE booking 
ALTER COLUMN payment_type SET DEFAULT 'Cash';

DROP TYPE payment_mode_enum;

COMMIT;
