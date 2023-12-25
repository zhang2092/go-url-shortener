DB_URL=postgresql://root:secret@localhost:5432/short_url?sslmode=disable

network:
	docker network create url-short-network

redis:
	docker run --name rd -d -p 6379:6379 redis:7.2.3 --requirepass "secret"

postgres:
	docker run --name postgres --network url-short-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root short_url

dropdb:
	docker exec -it postgres dropdb short_url

psql:
	docker exec -it postgres psql -U root -d short_url

migrateinit:
	migrate create -ext sql -dir db/schema -seq init_schema

migrateup:
	migrate -path db/schema -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/schema -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: network redis postgres createdb dropdb psql migrateup migratedown sqlc test server