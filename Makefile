
.PHONY: test proto

clean:
	@test ! -e bin || rm -r bin

build: clean
	go build 

test:
	@echo "Running test"
	go test -v ./...


run:
	@go run main.go

tidy:
	go mod tidy
