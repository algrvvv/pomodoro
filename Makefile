run:
	@go run cmd/pomodoro/main.go

web:
	@go run cmd/pomodoro/main.go --only-web

migration-create:
	@migrate create -ext sql -dir cmd/migrate/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migration-up:
	@go run cmd/migrate/main.go -up

migration-down:
	@go run cmd/migrate/main.go -down
