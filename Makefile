
build:
	mkdir -p ./bin
	go build -o ./bin/yubitoast ./src/cmd/yubitoast/main.go

run:
	go run ./src/cmd/yubitoast/main.go