.PHONY: build-producer build-consumer run-producer run-consumer up down logs clean
up:
	docker compose down -v && docker compose up --build