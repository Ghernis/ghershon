# Variables
BINARY_NAME=ghershon
CMD_PATH=./cmd/

# Detect OS for correct binary name
OS := $(shell uname -s)
ifeq ($(OS), Linux)
    EXT :=
    else
        EXT :=.exe
        endif
build:
	go build -o ghershon ./cmd/

build-win:
		GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe $(CMD_PATH)

run:
	./ghershon
clean:
	rm -f ghershon
