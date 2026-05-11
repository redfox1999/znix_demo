# 项目名称
BINARY_NAME=bin/server
SERVER_MAIN=cmd/server/main.go


.PHONY: all build run client clean help

# 默认执行编译
all: build

## build: 编译服务端
build:
	@mkdir -p bin
	@echo "Building server..."
	go build -o $(BINARY_NAME) $(SERVER_MAIN)

## run: 直接运行服务端
run:
	go run $(SERVER_MAIN)



## clean: 清理二进制文件
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

## help: 显示帮助信息
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
