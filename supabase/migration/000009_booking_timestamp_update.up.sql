BEGIN;

ALTER TABLE booking 
ALTER COLUMN booking_time TYPE TIMESTAMP WITH TIME ZONE,
ALTER COLUMN booking_time SET DEFAULT CURRENT_TIMESTAMP;

COMMIT;