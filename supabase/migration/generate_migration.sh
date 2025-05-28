#!/bin/bash

# Migration Analyzer Script
# Description: Consolidates all database migration files into single up and down files for analysis

MIGRATIONS_DIR="./"
OUTPUT_DIR="./migration_analysis"
UP_FILE="$OUTPUT_DIR/mega_up_migration.sql"
DOWN_FILE="$OUTPUT_DIR/mega_down_migration.sql"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Clear existing files
> "$UP_FILE"
> "$DOWN_FILE"

echo "=== Database Migration Analyzer ===" | tee -a "$UP_FILE" "$DOWN_FILE"
echo "Generated on: $(date)" | tee -a "$UP_FILE" "$DOWN_FILE"
echo "" | tee -a "$UP_FILE" "$DOWN_FILE"

# Function to add migration content with header
add_migration_content() {
    local file="$1"
    local output_file="$2"
    local migration_name=$(basename "$file")
    
    if [[ -f "$file" ]]; then
        echo "-- ========================================" >> "$output_file"
        echo "-- Migration: $migration_name" >> "$output_file"
        echo "-- ========================================" >> "$output_file"
        echo "" >> "$output_file"
        cat "$file" >> "$output_file"
        echo "" >> "$output_file"
        echo "" >> "$output_file"
    else
        echo "-- WARNING: $migration_name not found" >> "$output_file"
        echo "" >> "$output_file"
    fi
}

# Process UP migrations in ascending order
echo "Processing UP migrations..."
echo "-- CONSOLIDATED UP MIGRATIONS" >> "$UP_FILE"
echo "-- Apply in this order for fresh database setup" >> "$UP_FILE"
echo "" >> "$UP_FILE"

for file in $(ls "$MIGRATIONS_DIR"/*.up.sql 2>/dev/null | sort -V); do
    echo "Adding UP: $(basename "$file")"
    add_migration_content "$file" "$UP_FILE"
done

# Process DOWN migrations in descending order (reverse of up)
echo "Processing DOWN migrations..."
echo "-- CONSOLIDATED DOWN MIGRATIONS" >> "$DOWN_FILE"
echo "-- Apply in this order for complete rollback" >> "$DOWN_FILE"
echo "" >> "$DOWN_FILE"

for file in $(ls "$MIGRATIONS_DIR"/*.down.sql 2>/dev/null | sort -Vr); do
    echo "Adding DOWN: $(basename "$file")"
    add_migration_content "$file" "$DOWN_FILE"
done

# Generate summary
echo "" | tee -a "$UP_FILE" "$DOWN_FILE"
echo "-- MIGRATION SUMMARY" | tee -a "$UP_FILE" "$DOWN_FILE"
echo "-- Total UP migrations: $(ls "$MIGRATIONS_DIR"/*.up.sql 2>/dev/null | wc -l)" | tee -a "$UP_FILE" "$DOWN_FILE"
echo "-- Total DOWN migrations: $(ls "$MIGRATIONS_DIR"/*.down.sql 2>/dev/null | wc -l)" | tee -a "$UP_FILE" "$DOWN_FILE"

echo ""
echo "Migration analysis complete!"
echo "Files generated:"
echo "  - UP migrations: $UP_FILE"
echo "  - DOWN migrations: $DOWN_FILE"
echo ""
echo "Use these files to:"
echo "  1. Review complete database schema evolution"
echo "  2. Analyze migration dependencies"
echo "  3. Create fresh database setup script"
echo "  4. Plan rollback strategies"
