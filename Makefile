BINARY_NAME=telegram_translator

build:
	go build -ldflags="-s -w" -o $(BINARY_NAME) .

run:
	go run .

clean:
	rm -f $(BINARY_NAME)
