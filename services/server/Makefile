.PHONY: run-api

PORT := 8080
ADDR := "localhost"

run-api:
	go build -o streamserver && ./streamserver -port=${PORT} -addr=${ADDR}