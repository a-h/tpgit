build:
	GOOS=windows GOARCH=amd64 go build -o tpgit.exe
	GOOS=linux GOARCH=amd64 go build -o tpgit_linux
	GOOS=darwin GOARCH=amd64 go build -o tpgit_osx

test:
	go test ./... -cover