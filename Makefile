postgres:
	docker run --name postgres14 -p 5435:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres14 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5435/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5435/simple_bank?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v E:\Explore\GO\techschool\simplebank:/src -w /src kjconroy/sqlc generate

.PHONY: createdb dropdb postgres migrateup migratedown .sqlc