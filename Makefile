# Variables
GO=go
MIGRATION_DIR=./migration

# Commands
.PHONY: run-migrations rollback-migrations deploy seedData

run-migrations:
	@echo "Running migrations..."
	$(GO) run $(MIGRATION_DIR)/migration.go -direction up

rollback-migrations:
	@echo "Rolling back migrations..."
	$(GO) run $(MIGRATION_DIR)/migration.go -direction down

force-migration:
	@echo "Forcing migration version..."
	@read -p "Enter version to force (e.g., 0): " version; \
	$(GO) run $(MIGRATION_DIR)/migration.go -direction force -version $$version

deploy:
	@echo "Building and running the Go server..."
	$(GO) build -o skyfox-backend .
	./skyfox-backend

seedData:
	@echo "Seeding data..."
	./scripts/seed_data.sh
