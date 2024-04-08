build:
	go build -o ouro cmd/main.go

install: build
	mv ouro /usr/local/bin/
