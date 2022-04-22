postgres:
	docker run --name postgres-docker -p 8083:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14.2-alpine3.15
createdb:
	docker exec -it postgres-docker psql -U root -d root -c "CREATE DATABASE simple_bank;"

dropdb:
	docker exec -it postgres-docker psql -U root -d root -c "DROP DATABASE simple_bank;"

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:8083/simple_bank?sslmode=disable" -verbose up

migratedown:
	 migrate -path db/migration -database "postgresql://root:secret@localhost:8083/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test