PHONY: build run dev

# 开发模式
dev:
	cd internal/web && yarn dev & \
	DEV_MODE=true go run ./cmd/server

# 构建前端资源
build-frontend:
	cd internal/web && yarn build

# 生产构建
build: build-frontend
	go build -o app ./cmd/server

# 运行生产构建
run:
	./app