DB_SOURCE=postgresql://user:rocketman1@localhost:5432/peerbill_trader?sslmode=disable

start:
	sqlc init

generate:
	sqlc generate

init:
	docker run -it --rm --network host --volume "/Users/george/workspace/peerbill-trader-server/db:/db" migrate/migrate:v4.17.0 create -ext sql -dir /db/migrations $(name)

force:
	docker run -it --rm --network host --volume "/Users/george/workspace/peerbill-trader-server/db:/db" migrate/migrate:v4.17.0 create -ext sql -dir /db/migrations force $(version)

migrateup:
	docker run -it --rm --network host --volume ./db:/db migrate/migrate:v4.17.0 -path=/db/migrations -database "$(DB_SOURCE)" -verbose up

migrateup1:
	docker run -it --rm --network host --volume ./db:/db migrate/migrate:v4.17.0 -path=/db/migrations -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	docker run -it --rm --network host --volume ./db:/db migrate/migrate:v4.17.0 -path=/db/migrations -database "$(DB_SOURCE)" -verbose down

migratedown1:
	docker run -it --rm --network host --volume ./db:/db migrate/migrate:v4.17.0 -path=/db/migrations -database "$(DB_SOURCE)" -verbose down 1

test:
	go test -v -cover -short ./...

mock:
	mockgen -package mockdb -destination db/mock/repository.go peerbill-trader-server/db/sqlc DatabaseContract

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto


.PHONY: start generate init migrateup migratedown mock test migrateup1 migratedown1 force proto
