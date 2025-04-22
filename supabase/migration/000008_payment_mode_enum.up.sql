BEGIN;

CREATE TYPE payment_mode_enum AS ENUM ('Cash', 'Card');

ALTER TABLE booking ALTER COLUMN payment_type DROP DEFAULT;

ALTER TABLE booking 
ALTER COLUMN payment_type TYPE payment_mode_enum 
USING 
    CASE 
        WHEN payment_type = 'Cash' THEN 'Cash'::payment_mode_enum
        WHEN payment_type = 'Card' THEN 'Card'::payment_mode_enum
        ELSE 'Cash'::payment_mode_enum
    END;

ALTER TABLE booking 
ALTER COLUMN payment_type SET DEFAULT 'Cash'::payment_mode_enum;

COMMIT;
