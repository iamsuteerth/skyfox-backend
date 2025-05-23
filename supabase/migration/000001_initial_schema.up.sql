BEGIN;

CREATE TYPE booking_status AS ENUM ('Pending', 'Confirmed', 'CheckedIn');
CREATE TYPE user_role_enum AS ENUM ('admin', 'customer', 'staff');
CREATE TABLE usertable (
  id BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  username VARCHAR(30) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  role user_role_enum NOT NULL
);
CREATE TABLE password_history (
  id BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  username VARCHAR(30) UNIQUE NOT NULL,
  previous_password_1 VARCHAR(255),
  previous_password_2 VARCHAR(255),
  previous_password_3 VARCHAR(255),
  CONSTRAINT fk_password_history_username FOREIGN KEY (username) REFERENCES usertable(username) ON DELETE CASCADE
);
CREATE TABLE stafftable (
  id BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  username VARCHAR(30) NOT NULL,
  name VARCHAR(70) NOT NULL,
  counter_no INTEGER,
  CONSTRAINT fk_stafftable_username FOREIGN KEY (username) REFERENCES usertable(username) ON DELETE CASCADE
);
CREATE TABLE customertable (
  id BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  name VARCHAR(70) NOT NULL,
  username VARCHAR(30) UNIQUE NOT NULL,
  number VARCHAR(10) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  profile_img VARCHAR(255),
  CONSTRAINT fk_customertable_username FOREIGN KEY (username) REFERENCES usertable(username) ON DELETE CASCADE,
  CONSTRAINT check_number CHECK (number ~ '^[0-9]{10}$')
);
CREATE TABLE admin_booked_customer (
  id BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  name VARCHAR(70) NOT NULL,
  number VARCHAR(10) NOT NULL,
  CONSTRAINT check_number CHECK (number ~ '^[0-9]{10}$')
);
CREATE TABLE seat (
  seat_number VARCHAR(10) PRIMARY KEY NOT NULL,
  seat_type VARCHAR(50)
);
CREATE TABLE slot (
  id BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  name VARCHAR(50) NOT NULL,
  start_time TIME NOT NULL,
  end_time TIME NOT NULL
);
CREATE TABLE show (
  id BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  movie_id VARCHAR(30) NOT NULL,
  date DATE NOT NULL,
  slot_id BIGINT NOT NULL,
  cost DECIMAL(10, 2) NOT NULL,
  CONSTRAINT fk_show_slot FOREIGN KEY (slot_id) REFERENCES slot(id),
  CONSTRAINT unique_slot_date UNIQUE (slot_id, date)
);
CREATE TABLE booking (
  id BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  date DATE NOT NULL,
  show_id BIGINT,
  customer_id BIGINT,
  no_of_seats INTEGER NOT NULL,
  amount_paid DECIMAL(10, 2) NOT NULL,
  status booking_status NOT NULL DEFAULT 'Pending',
  booking_time TIMESTAMP NOT NULL DEFAULT NOW(),
  payment_type VARCHAR(50) NOT NULL DEFAULT 'Cash',
  CONSTRAINT fk_booking_show FOREIGN KEY (show_id) REFERENCES show(id),
  CONSTRAINT fk_booking_customer FOREIGN KEY (customer_id) REFERENCES admin_booked_customer(id)
);
CREATE TABLE booking_seat_mapping (
  id BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  booking_id BIGINT NOT NULL,
  seat_number VARCHAR(10) NOT NULL,
  CONSTRAINT fk_booking_seat_mapping_booking FOREIGN KEY (booking_id) REFERENCES booking(id) ON DELETE CASCADE,
  CONSTRAINT fk_booking_seat_mapping_seat FOREIGN KEY (seat_number) REFERENCES seat(seat_number)
);

COMMIT;