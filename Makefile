swag:
	swag init -g ./internal/app/app.go
test:
	go test -v -cover ./internal/repository/pgdb  
	go test -v -cover ./internal/controller/http/v1
migrate:
	migrate -path ./schema -database 'postgres://postgres:postgres@0.0.0.0:5432/postgres?sslmode=disable' up
build:
	docker-compose build
run:
	docker-compose up
create-m:
	migrate create -ext sql -dir ./schema -seq init