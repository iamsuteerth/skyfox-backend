#!/bin/bash

# Supabase connection details
DB_HOST="${DB_HOST}"
DB_NAME="postgres"
DB_USER="postgres"
DB_PASSWORD="${DB_PASSWORD}"
DB_PORT="5432"

# Movie IDs list
MOVIE_IDS="tt6644200 tt6857112 tt7784604 tt5052448 tt1396484 tt5968394 tt4972582 tt6823368 tt7556122 tt1179933 \
           tt7349950 tt2935510 tt0437086 tt4154664 tt3016748 tt3513498 tt4178092 tt4154796 tt2382320 tt1833116"

# Functions
get_random_movie_id() {
  movie_ids_array=($MOVIE_IDS)
  random_index=$((RANDOM % ${#movie_ids_array[@]}))
  echo "${movie_ids_array[$random_index]}"
}

get_random_price() {
  price_lower_value=150
  price_upper_value=300
  price=$((RANDOM % (price_upper_value - price_lower_value + 1) + price_lower_value))
  echo "$price.$((RANDOM % 99))"
}

# Clear existing data
clear_existing_data() {
  echo "Clearing existing data..."
  PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -q -c "
    -- Disable triggers temporarily to avoid foreign key issues
    SET session_replication_role = 'replica';
    
    -- Truncate tables in the correct order
    TRUNCATE booking_seat_mapping, booking, show, seat, slot CASCADE;
    
    -- Reset sequences
    ALTER SEQUENCE slot_id_seq RESTART WITH 1;
    ALTER SEQUENCE show_id_seq RESTART WITH 1;
    ALTER SEQUENCE booking_id_seq RESTART WITH 1;
    ALTER SEQUENCE booking_seat_mapping_id_seq RESTART WITH 1;
    
    -- Re-enable triggers
    SET session_replication_role = 'origin';
  "
  echo "✓ Data cleared successfully and sequences reset"
}

# Seed slot data
seed_slot_data() {
  echo "Seeding slot data..."
  PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -q -c "
  INSERT INTO slot (name, start_time, end_time) VALUES 
    ('Morning', '08:30:00', '11:30:00'),
    ('Afternoon', '12:30:00', '15:30:00'),
    ('Evening', '16:30:00', '19:30:00'),
    ('Night', '20:30:00', '23:30:00');"
  
  slot_count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM slot;")
  echo "✓ Added $slot_count time slots"
}

# Seed seat data
seed_seat_data() {
  echo "Seeding seat data..."
  
  # Create a temporary file with all INSERT statements
  temp_file=$(mktemp)
  
  echo "BEGIN;" > $temp_file
  # Rows A-E are Standard
  for row in {A..E}; do
    for num in {1..10}; do
      seat="${row}${num}"
      echo "INSERT INTO seat (seat_number, seat_type) VALUES ('$seat', 'Standard');" >> $temp_file
    done
  done
  
  # Rows F-J are Deluxe
  for row in {F..J}; do
    for num in {1..10}; do
      seat="${row}${num}"
      echo "INSERT INTO seat (seat_number, seat_type) VALUES ('$seat', 'Deluxe');" >> $temp_file
    done
  done
  echo "COMMIT;" >> $temp_file
  
  # Execute all inserts in a single transaction (quietly)
  PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -q -f $temp_file
  
  # Count and display how many seats were added
  standard_count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM seat WHERE seat_type = 'Standard';")
  deluxe_count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM seat WHERE seat_type = 'Deluxe';")
  
  echo "✓ Added $(echo $standard_count | xargs) Standard seats (Rows A-E)"
  echo "✓ Added $(echo $deluxe_count | xargs) Deluxe seats (Rows F-J)"
  
  # Remove temporary file
  rm $temp_file
}

# Seed show data
seed_show_data() {
  echo "Seeding show data for the next 21 days..."
  
  # Get today's date
  today=$(date +%Y-%m-%d)
  
  # Get slot IDs
  slot_ids=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT id FROM slot ORDER BY id;")
  
  # Create a temporary file with all INSERT statements
  temp_file=$(mktemp)
  
  echo "BEGIN;" > $temp_file
  
  # For each day in the next 21 days
  for i in {0..20}; do
    # Calculate date
    if [[ "$OSTYPE" == "darwin"* ]]; then
      # macOS
      current_date=$(date -j -v+${i}d -f "%Y-%m-%d" "$today" +%Y-%m-%d)
    else
      # Linux
      current_date=$(date -d "$today + $i days" +%Y-%m-%d)
    fi
    
    # For each slot
    for slot_id in $slot_ids; do
      movie_id=$(get_random_movie_id)
      price=$(get_random_price)
      
      echo "INSERT INTO show (movie_id, date, slot_id, cost) VALUES ('$movie_id', '$current_date', $slot_id, $price);" >> $temp_file
    done
  done
  
  echo "COMMIT;" >> $temp_file
  
  # Execute all inserts in a single transaction (quietly)
  PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -q -f $temp_file
  
  # Count how many shows were added
  show_count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM show;")
  days_count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(DISTINCT date) FROM show;")
  
  echo "✓ Added $(echo $show_count | xargs) shows across $(echo $days_count | xargs) days"
  
  # Remove temporary file
  rm $temp_file
}

# Main execution
echo "Starting data seeding process..."
clear_existing_data
seed_slot_data
seed_seat_data
seed_show_data
echo "✓ Data seeding complete!"