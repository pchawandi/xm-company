
build-docker:
	docker compose build --no-cache

run-local:
	docker start dockerPostgres
	export POSTGRES_DB=go_app_dev
	export POSTGRES_USER=docker
	export POSTGRES_PASSWORD=password
	export POSTGRES_PORT=5435
	export JWT_SECRET_KEY=ObL89O3nOSSEj6tbdHako0cXtPErzBUfq8l8o/3KD9g=INSECURE
	export API_SECRET_KEY=cJGZ8L1sDcPezjOy1zacPJZxzZxrPObm2Ggs1U0V+fE=INSECURE
	export POSTGRES_HOST=localhost
	go run cmd/server/main.go

up:
	docker compose up --build

down:
	docker compose down

restart:
	docker compose restart

build:
	go build -v ./...

test:
	go test -v ./... -race -cover

clean:
	docker stop xm-company
	docker stop dockerPostgres
	docker rm xm-company
	docker rm dockerPostgres
	docker image rm xm-company-backend
	rm -rf .dbdata
