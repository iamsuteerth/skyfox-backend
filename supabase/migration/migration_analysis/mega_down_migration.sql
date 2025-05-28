=== Database Migration Analyzer ===
Generated on: Wed May 28 06:41:59 PM IST 2025

-- CONSOLIDATED DOWN MIGRATIONS
-- Apply in this order for complete rollback

-- ========================================
-- Migration: 000021_optimize_database_indices.down.sql
-- ========================================

BEGIN;

DROP INDEX IF EXISTS idx_booking_date_status;
DROP INDEX IF EXISTS idx_show_date_slot;
DROP INDEX IF EXISTS idx_wallet_transaction_user_time;
DROP INDEX IF EXISTS idx_payment_transaction_booking_status;
DROP INDEX IF EXISTS idx_active_bookings;
DROP INDEX IF EXISTS idx_active_reset_tokens;
DROP INDEX IF EXISTS idx_wallet_balance_lookup;
DROP INDEX IF EXISTS idx_wallet_txn_aggregation;
DROP INDEX IF EXISTS idx_booking_details_covering;
DROP INDEX IF EXISTS idx_show_movie_date;

CREATE INDEX idx_booking_seat_mapping_seat_number ON booking_seat_mapping(seat_number);
CREATE INDEX idx_stafftable_username ON stafftable(username);

COMMIT;


-- ========================================
-- Migration: 000020_wallet_customer_cascade.down.sql
-- ========================================

BEGIN;

ALTER TABLE wallet_transaction DROP CONSTRAINT fk_wallet_username;
ALTER TABLE customer_wallet DROP CONSTRAINT fk_customer_wallet_username;

ALTER TABLE wallet_transaction 
ADD CONSTRAINT fk_wallet_username 
FOREIGN KEY (username) REFERENCES customertable(username);

ALTER TABLE customer_wallet 
ADD CONSTRAINT fk_customer_wallet_username 
FOREIGN KEY (username) REFERENCES customertable(username);

COMMIT;


-- ========================================
-- Migration: 000019_wallet_transaction_table.down.sql
-- ========================================

BEGIN;

DROP TABLE IF EXISTS wallet_transaction;

COMMIT;

-- ========================================
-- Migration: 000018_wallet_transaction_type_enum.down.sql
-- ========================================

BEGIN;

DROP TYPE IF EXISTS wallet_transaction_type;

COMMIT;

-- ========================================
-- Migration: 000017_wallet_db_tables.down.sql
-- ========================================

BEGIN;

DROP INDEX IF EXISTS idx_wallet_username;

DROP TABLE IF EXISTS customer_wallet;

COMMIT;

-- ========================================
-- Migration: 000016_wallet_payment_type.down.sql
-- ========================================

BEGIN;

COMMIT;

-- ========================================
-- Migration: 000015_performance_idices.down.sql
-- ========================================

BEGIN;

DROP INDEX IF EXISTS idx_booking_status;
DROP INDEX IF EXISTS idx_booking_booking_time;
DROP INDEX IF EXISTS idx_show_slot_id;
DROP INDEX IF EXISTS idx_show_movie_id;
DROP INDEX IF EXISTS idx_admin_booked_customer_booking_id;

COMMIT;

-- ========================================
-- Migration: 000014_circular_dependency_fix.down.sql
-- ========================================

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

-- ========================================
-- Migration: 000013_add_foreign_key_indices.down.sql
-- ========================================

BEGIN;

DROP INDEX IF EXISTS idx_booking_customer_username;
DROP INDEX IF EXISTS idx_booking_show_id;
DROP INDEX IF EXISTS idx_booking_seat_mapping_booking_id;
DROP INDEX IF EXISTS idx_booking_seat_mapping_seat_number;
DROP INDEX IF EXISTS idx_stafftable_username;

COMMIT;

-- ========================================
-- Migration: 000012_pending_booking_tracker.down.sql
-- ========================================

BEGIN;

DROP TABLE IF EXISTS pending_booking_tracker;

COMMIT;

-- ========================================
-- Migration: 000011_payment_txd_table.down.sql
-- ========================================

BEGIN;

DROP TABLE IF EXISTS payment_transaction;

COMMIT;

-- ========================================
-- Migration: 000010_relationship_btw_abc_booking.down.sql
-- ========================================

BEGIN;

ALTER TABLE admin_booked_customer DROP CONSTRAINT IF EXISTS fk_admin_customer_booking;

ALTER TABLE booking DROP CONSTRAINT IF EXISTS unique_customer_id; 

ALTER TABLE booking
ADD CONSTRAINT fk_booking_customer
FOREIGN KEY (customer_id) 
REFERENCES admin_booked_customer(id)
ON DELETE NO ACTION;

COMMIT;

-- ========================================
-- Migration: 000009_booking_timestamp_update.down.sql
-- ========================================

BEGIN;

ALTER TABLE booking 
ALTER COLUMN booking_time TYPE TIMESTAMP,
ALTER COLUMN booking_time SET DEFAULT NOW();

COMMIT;


-- ========================================
-- Migration: 000008_payment_mode_enum.down.sql
-- ========================================

BEGIN;

ALTER TABLE booking ALTER COLUMN payment_type DROP DEFAULT;

ALTER TABLE booking 
ALTER COLUMN payment_type TYPE VARCHAR(50) 
USING payment_type::VARCHAR(50);

ALTER TABLE booking 
ALTER COLUMN payment_type SET DEFAULT 'Cash';

DROP TYPE payment_mode_enum;

COMMIT;


-- ========================================
-- Migration: 000007_booking_table_modifications.down.sql
-- ========================================

BEGIN;

ALTER TABLE booking DROP CONSTRAINT check_customer_type;
ALTER TABLE booking DROP CONSTRAINT fk_booking_customer_username;
ALTER TABLE booking DROP COLUMN customer_username;

COMMIT;


-- ========================================
-- Migration: 000006_created_at_column.down.sql
-- ========================================

BEGIN;

ALTER TABLE usertable
DROP COLUMN created_at;

COMMIT;

-- ========================================
-- Migration: 000005_add_unique_constraint_email_token.down.sql
-- ========================================

BEGIN;

ALTER TABLE password_reset_tokens 
DROP CONSTRAINT unique_email_token;

COMMIT;


-- ========================================
-- Migration: 000004_security_question_temp_tkn.down.sql
-- ========================================

BEGIN;

DROP INDEX IF EXISTS idx_password_reset_tokens_token;
DROP INDEX IF EXISTS idx_password_reset_tokens_email;
DROP TABLE IF EXISTS password_reset_tokens;

COMMIT;

-- ========================================
-- Migration: 000003_add_security_questions.down.sql
-- ========================================

BEGIN;

DROP INDEX IF EXISTS idx_security_question_id;

ALTER TABLE customertable 
DROP COLUMN IF EXISTS security_question_id,
DROP COLUMN IF EXISTS security_answer_hash;

DROP TABLE IF EXISTS security_questions;

COMMIT;

-- ========================================
-- Migration: 000002_seat_types.down.sql
-- ========================================

BEGIN;

DROP TYPE IF EXISTS seat_type_enum;

COMMIT;

-- ========================================
-- Migration: 000001_initial_schema.down.sql
-- ========================================

BEGIN;

DROP TABLE IF EXISTS booking_seat_mapping;
DROP TABLE IF EXISTS booking;                    
DROP TABLE IF EXISTS admin_booked_customer;     
DROP TABLE IF EXISTS show;
DROP TABLE IF EXISTS slot;
DROP TABLE IF EXISTS seat;
DROP TABLE IF EXISTS customertable;
DROP TABLE IF EXISTS stafftable;
DROP TABLE IF EXISTS password_history;
DROP TABLE IF EXISTS usertable;
DROP TYPE IF EXISTS booking_status;
DROP TYPE IF EXISTS user_role_enum;

COMMIT;


-- MIGRATION SUMMARY
-- Total UP migrations: 21
-- Total DOWN migrations: 21
