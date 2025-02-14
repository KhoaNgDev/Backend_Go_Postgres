postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine || true

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank || echo "Database already exists!"

dropdb:
	docker exec -it postgres12 dropdb simple_bank || echo "Database does not exist!"

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up || echo "Migration failed!"

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down || echo "Migration rollback failed!"

checkmigrate: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" version

forcefix:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" force 1 || echo "Force migration failed!"

sqlc:
	sqlc generate	

test:
	go test -v -cover ./...
.PHONY: postgres createdb dropdb migrateup migratedown sqlc checkmigrate forcefix test
