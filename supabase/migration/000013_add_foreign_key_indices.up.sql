BEGIN;

CREATE INDEX idx_booking_customer_username ON booking(customer_username);

CREATE INDEX idx_booking_show_id ON booking(show_id);

CREATE INDEX idx_booking_seat_mapping_booking_id ON booking_seat_mapping(booking_id);

CREATE INDEX idx_booking_seat_mapping_seat_number ON booking_seat_mapping(seat_number);

CREATE INDEX idx_stafftable_username ON stafftable(username);

COMMIT;
