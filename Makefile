.PHONY: build run clean install

build:
	go build -o clashtui .

run:
	go run .

clean:
	rm -f clashtui

install:
	go build -o clashtui .
	cp clashtui ~/.local/bin/clashtui
	@echo "Installed to ~/.local/bin/clashtui"