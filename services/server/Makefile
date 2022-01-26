.PHONY: run-api run-watch


run-watch:
	npx nodemon --watch './**/*.go' --signal SIGTERM --exec 'go run main.go'

run-api:
	go build -o streamserver && \
 	export $$(grep -v '^#' .env | xargs) && \
	./streamserver

test:
	go test ./app/... -v

test-color:
	go test ./app/... -v | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''