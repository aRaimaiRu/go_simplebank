DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

network:
	docker network create bank-network
postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank
migratecreate:
	migrate create -ext sql -dir db/migration -seq init_schema
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
dropdb:
	docker exec -it postgres dropdb simple_bank
test:
	go test -v -cover ./...
sqlc:
	sqlc generate
mock:
	mockgen -package mock_db -destination  db/mock/store.go go_simplebank/db/sqlc Store
server:
	go run main.go