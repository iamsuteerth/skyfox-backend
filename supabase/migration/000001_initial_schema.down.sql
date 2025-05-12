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