.PHONY: build build-win dev run clean

PROJECT="game-gorl"
BUILD_PATH="./build"

init:
	mkdir build
	mkdir runtime

build-win:
	CGO_LDFLAGS="-static-libgcc -static -lpthread"\
				GOOS=windows GOARCH=amd64 CGO_ENABLED=1\
				CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++\
				go build -a -x -o $(BUILD_PATH)/$(PROJECT)-win.exe

build:
	mkdir -p $(BUILD_PATH)
	cp -r assets/* $(BUILD_PATH)
	go build -o $(BUILD_PATH)/$(PROJECT)-linux -v

build-debug:
	mkdir -p $(BUILD_PATH)
	cp -r assets/* $(BUILD_PATH)
	CGO_CFLAGS='-O0 -g' go build -a -v -gcflags="all=-N -l" -o $(BUILD_PATH)/$(PROJECT)-linux

run:
	cd $(BUILD_PATH); ./$(PROJECT)-linux

dev:
	@make build && make run || echo "build failed!"

clean:
	rm -r $(BUILD_PATH)/*
	mkdir -p $(BUILD_PATH)
	cp -r assets/* $(BUILD_PATH)

lint:
	nilaway main.go
