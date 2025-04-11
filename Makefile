run:
	docker compose up -d
	go run cmd/numeral/main.go

.PHONY: run
