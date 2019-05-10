release:
	mkdir -p release/macos
	mkdir release/linux
	mkdir release/windows
	GOOS=darwin go build -o release/macos/multiwatch
	GOOS=linux go build -o release/linux/multiwatch
	GOOS=windows go build -o release/windows/multiwatch.exe
