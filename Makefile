BINARY_NAME=telegram_translator

build:
	go build -ldflags="-s -w" -o $(BINARY_NAME) main.go

run:
	go run main.go

clean:
	rm -f $(BINARY_NAME)
