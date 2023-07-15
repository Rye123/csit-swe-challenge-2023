.DEFAULT_GOAL := build
MAIN="./cmd/main/server.go"
TEST_DB="./internal/db/"

build:
	go build ${MAIN}

run:
	go run ${MAIN}

test:
	go test ${TEST_DB} -v
