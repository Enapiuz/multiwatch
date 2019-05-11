build:
	mkdir -p release/macos
	mkdir -p release/linux
	mkdir -p release/windows
	GOOS=darwin go build -o release/macos/multiwatch
	zip release/macos.zip release/macos/multiwatch
	GOOS=linux go build -o release/linux/multiwatch
	zip release/linux.zip release/linux/multiwatch
	GOOS=windows go build -o release/windows/multiwatch.exe
	zip release/windows.zip release/windows/multiwatch.exe

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...
