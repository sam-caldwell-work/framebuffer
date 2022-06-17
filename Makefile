help:
	@echo "'make test' will run unit tests."

.PHONY: test
test:
	@go test -v
