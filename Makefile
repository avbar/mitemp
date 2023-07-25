build:
	docker compose build

up:
	docker compose up

run-all:
	make build
	make up