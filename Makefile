run:
	pbpaste | go run main.go

build:
	go build -o build/log-beautify main.go

install:
	go install github.com/mylxsw/log-beautify@latest

clean:
	rm -rf build

.PHONY: run build install clean
