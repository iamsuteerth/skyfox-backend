BEGIN;

CREATE TABLE pending_booking_tracker (
    booking_id BIGINT PRIMARY KEY,
    expiration_time TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT fk_pending_booking_tracker_booking FOREIGN KEY (booking_id) REFERENCES booking(id) ON DELETE CASCADE
);

COMMIT;
