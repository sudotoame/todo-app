include .env
export

run:
	@go run cmd/todo-app/main.go

migrate-version:
	@migrate -path migrations -database ${PGX_CONN_LOCAL} version

migrate-up:
	@migrate -path migrations -database ${PGX_CONN_LOCAL} up

migrate-down:
	@migrate -path migrations -database ${PGX_CONN_LOCAL} down