DB_SOURCE=postgresql://user:rocketman1@localhost:5432/peerbill_trader?sslmode=disable

start:
	sqlc init

generate:
	sqlc generate

init:
	docker run -it --rm --network host --volume "/Users/george/workspace/peerbill-trader-server/db:/db" migrate/migrate:v4.17.0 create -ext sql -dir /db/migrations init_schema

migrateup:
	docker run -it --rm --network host --volume ./db:/db migrate/migrate:v4.17.0 -path=/db/migrations -database "$(DB_SOURCE)" -verbose up

migratedown:
	docker run -it --rm --network host --volume ./db:/db migrate/migrate:v4.17.0 -path=/db/migrations -database "$(DB_SOURCE)" -verbose down

.PHONY: start generate init migrateup migratedown
