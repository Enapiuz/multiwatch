GOPACKAGES=$(shell find . -name '*.go' -not -path "./vendor/*" -exec dirname {} \; | uniq)
GOFILES_NOVENDOR=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

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

watch:
	go run main.go

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic -v $(GOPACKAGES)

vet:
	go vet $(GOPACKAGES)

lint:
	ls $(GOFILES_NOVENDOR) | xargs -L1 golint -set_exit_status

imports:
	goimports -w $(GOFILES_NOVENDOR)
