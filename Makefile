PHONY: build run dev

# 开发模式
dev:
	cd ./web && yarn dev & \
	DEV_MODE=true go run ./cmd/server

# 构建前端资源
build-web:
	cd ./web && yarn build

# 生产构建
build: build-web
	go build -o yop ./cmd/server

# 运行生产构建
run:
	./yop