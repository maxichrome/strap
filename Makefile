build:
	go build -o bin/strap

run: build
	@./bin/strap
