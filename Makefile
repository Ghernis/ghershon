# Variables
BINARY_NAME=gher
CMD_PATH=./cmd/

# Detect OS for correct binary name
OS := $(shell uname -s)
ifeq ($(OS), Linux)
    EXT :=
    else
        EXT :=.exe
        endif
build:
	go build -o $(BINARY_NAME) ./cmd/

build-win:
		GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe $(CMD_PATH)
rp:
	cp $(BINARY_NAME) ~/bin
	chmod 744 ~/bin/$(BINARY_NAME)

run:
	./$(BINARY_NAME)
r:
	go run ./cmd/main.go
clean:
	rm -f $(BINARY_NAME)
