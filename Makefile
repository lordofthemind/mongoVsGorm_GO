# Environment variables for MongoDB
MONGODB_CONTAINER_NAME ?= MongoVsGorm_MG
MONGODB_IMAGE_TAG ?= latest
MONGODB_DB_NAME ?= MongoVsGorm_MGDB
MONGODB_PORT ?= 27017

# Environment variables for PostgreSQL
PG_CONTAINER_NAME ?= MongoVsGorm_PG
PG_IMAGE_TAG ?= latest
PG_DB_NAME ?= MongoVsGorm_PGDB
PG_DB_USERNAME ?= postgres
PG_DB_PASSWORD ?= MongoVsGormSecret
PG_PORT ?= 5432

# Colors for help command
CYAN := \033[36m
RESET := \033[0m

# Docker MongoDB commands
crtmgdb: ## Create and start the MongoDB container
	@echo "Creating and starting MongoDB container..."
	docker run --name $(MONGODB_CONTAINER_NAME) -p $(MONGODB_PORT):27017 -d mongo:$(MONGODB_IMAGE_TAG)

strmgdb: ## Start the MongoDB container
	@echo "Starting MongoDB container..."
	docker start $(MONGODB_CONTAINER_NAME)

stpmgdb: ## Stop the MongoDB container
	@echo "Stopping MongoDB container..."
	docker stop $(MONGODB_CONTAINER_NAME)

rmvmgdb: ## Remove the MongoDB container
	@echo "Removing MongoDB container..."
	docker rm $(MONGODB_CONTAINER_NAME)

# MongoDB database commands
createdb_mongodb: strmgdb ## Create MongoDB database
	@echo "Creating MongoDB database..."
	docker exec -it $(MONGODB_CONTAINER_NAME) mongosh --eval "use $(MONGODB_DB_NAME)"

dropdb_mongodb: strmgdb ## Drop MongoDB database
	@echo "Dropping MongoDB database..."
	docker exec -it $(MONGODB_CONTAINER_NAME) mongosh --eval "db.getSiblingDB('$(MONGODB_DB_NAME)').dropDatabase()"

# Docker PostgreSQL commands
crtpgdb: ## Create and start the PostgreSQL container
	@echo "Creating and starting PostgreSQL container..."
	docker run --name $(PG_CONTAINER_NAME) -p $(PG_PORT):5432 -e POSTGRES_DB=$(PG_DB_NAME) -e POSTGRES_USER=$(PG_DB_USERNAME) -e POSTGRES_PASSWORD=$(PG_DB_PASSWORD) -d postgres:$(PG_IMAGE_TAG)

strpgdb: ## Start the PostgreSQL container
	@echo "Starting PostgreSQL container..."
	docker start $(PG_CONTAINER_NAME)

stppgdb: ## Stop the PostgreSQL container
	@echo "Stopping PostgreSQL container..."
	docker stop $(PG_CONTAINER_NAME)

rmvpgdb: ## Remove the PostgreSQL container
	@echo "Removing PostgreSQL container..."
	docker rm $(PG_CONTAINER_NAME)

# PostgreSQL database commands
createdb_pgdb: strpgdb ## Create PostgreSQL database
	@echo "Creating PostgreSQL database..."
	docker exec -it $(PG_CONTAINER_NAME) createdb -U $(PG_DB_USERNAME) $(PG_DB_NAME)

dropdb_pgdb: strpgdb ## Drop PostgreSQL database
	@echo "Dropping PostgreSQL database..."
	docker exec -it $(PG_CONTAINER_NAME) dropdb -U $(PG_DB_USERNAME) $(PG_DB_NAME)

# Help command
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1m%-12s\033[0m %s\n\n", "Command", "Description"} /^[a-zA-Z_-]+:.*?##/ { printf "\033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: crtmgdb strmgdb stpmgdb rmvmgdb createdb_mongodb dropdb_mongodb crtpgdb strpgdb stppgdb rmvpgdb createdb_pgdb dropdb_pgdb help
