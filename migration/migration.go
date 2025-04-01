package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	migrationsPath = "./supabase/migration"
	dbName         = "postgres"
)

func main() {
	var direction string
	var forceVersion string
	flag.StringVar(&direction, "direction", "up", "Migration direction (up, down, or force)")
	flag.StringVar(&forceVersion, "version", "", "Version to force (only for force command)")
	flag.Parse()

	dbHost := getEnv("DB_HOST", "")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("POSTGRES_USERNAME", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")

	if dbPassword == "" {
		fmt.Println("Error: DB_PASSWORD environment variable is not set")
		os.Exit(1)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("Error pinging database: %v\n", err)
		os.Exit(1)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Printf("Error creating migration driver: %v\n", err)
		os.Exit(1)
	}

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		fmt.Printf("Error getting absolute path: %v\n", err)
		os.Exit(1)
	}
	sourcePath := fmt.Sprintf("file://%s", absPath)

	m, err := migrate.NewWithDatabaseInstance(sourcePath, dbName, driver)
	if err != nil {
		fmt.Printf("Error creating migrator: %v\n", err)
		os.Exit(1)
	}

	switch direction {
	case "up":
		fmt.Println("Running migrations up...")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("Error running migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		fmt.Println("Rolling back migrations...")
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("Error rolling back migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Rollback completed successfully")

	case "force":
		if forceVersion == "" {
			fmt.Println("Error: version is required for force command")
			os.Exit(1)
		}

		version, err := strconv.ParseUint(forceVersion, 10, 64)
		if err != nil {
			fmt.Printf("Error parsing version: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Forcing migration version to %d...\n", version)
		if err := m.Force(int(version)); err != nil {
			fmt.Printf("Error forcing migration version: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully forced version to %d\n", version)

	default:
		fmt.Printf("Unknown direction: %s (use 'up', 'down', or 'force')\n", direction)
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
