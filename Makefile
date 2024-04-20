# Docker container name
CONTAINER_NAME = backendMasterGo
# PostgreSQL database name
DB_NAME = simple_bank

# Create PostgreSQL container
createpg:
	docker run --name $(CONTAINER_NAME) -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=backendMasterGoSecret -d postgres:16-alpine

# Start PostgreSQL container
startpg:
	docker start $(CONTAINER_NAME)

# Stop PostgreSQL container
stoppg:
	docker stop $(CONTAINER_NAME)

# Remove PostgreSQL container
removepg:
	docker rm $(CONTAINER_NAME)

# Connect to PostgreSQL database using psql
psql:
	docker exec -it $(CONTAINER_NAME) psql -U root $(DB_NAME)

# Open shell in PostgreSQL container
sh:
	docker exec -it $(CONTAINER_NAME) /bin/sh

# Create new database
createdb:
	docker exec -it $(CONTAINER_NAME) createdb --username=root --owner=root $(DB_NAME)

# Drop specified database
dropdb:
	docker exec -it $(CONTAINER_NAME) dropdb $(DB_NAME)

# Dump database schema and data into SQL file
dumpdb:
	docker exec -it $(CONTAINER_NAME) pg_dump -U root -d $(DB_NAME) > dump.sql

# Restore database from SQL dump file
restoredb:
	docker exec -i $(CONTAINER_NAME) psql -U root -d $(DB_NAME) < dump.sql

# Connect to database using psql
connectdb:
	docker exec -it $(CONTAINER_NAME) psql -U root -d $(DB_NAME)

migrateup:
	migrate -path db/migration -database "postgresql://root:backendMasterGoSecret@localhost:5432/$(DB_NAME)?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:backendMasterGoSecret@localhost:5432/$(DB_NAME)?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:backendMasterGoSecret@localhost:5432/$(DB_NAME)?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:backendMasterGoSecret@localhost:5432/$(DB_NAME)?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/lordofthemind/backendMasterGo/db/sqlc Store

# Phony targets to avoid conflicts with files of the same name
.PHONY: createpg startpg stoppg removepg psql sh createdb dropdb dumpdb restoredb connectdb migrateup migratedown sqlc test server mock migrateup1 migratedown1
