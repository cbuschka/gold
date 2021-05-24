TOP_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

test:
	@echo "### Running unit tests..."; \
	cd ${TOP_DIR}
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic ./internal/... ./cmd/...

tidy_mod:
	@echo "### Updating deps..."; \
	cd ${TOP_DIR}
	go mod tidy

run_daemon:
	@echo "### Running golfd..."; \
	cd ${TOP_DIR}
	go run ./cmd/golfd/golfd.go

run_query_list:
	@echo "### Running golfq..."; \
	cd ${TOP_DIR}
	go run ./cmd/golfq/golfq.go list
