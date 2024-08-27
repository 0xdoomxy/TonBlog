# 定义变量
DOCKER_COMPOSE_FILE = deployment/docker-compose.yaml
DOCKER_IMAGE = backend:latest
BACKEND_DIR = backend
FRONTEND_DIR = frontend

# 默认目标
all: build

# 检查并停止正在运行的 Docker Compose 服务
check_and_stop:
	@if docker-compose -f $(DOCKER_COMPOSE_FILE) ps -q | grep -q .; then \
		echo "Stopping running Docker Compose services..."; \
		docker-compose -f $(DOCKER_COMPOSE_FILE) stop; \
	fi

# 删除 Docker 镜像
remove_docker_image:
	@if docker images -q $(DOCKER_IMAGE) | grep -q .; then \
		echo "Removing Docker image $(DOCKER_IMAGE)..."; \
		docker rmi $(DOCKER_IMAGE); \
	fi

# 编译后端程序
build_backend:
	cd $(BACKEND_DIR) && go build -ldflags="-s -w" main.go && chmod u+x blog

# 打包前端
build_frontend:
	cd $(FRONTEND_DIR) && npm install && npm run build

# 运行 Docker Compose
up_docker_compose:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# 主目标
build: check_and_stop remove_docker_image build_backend build_frontend 

run: up_docker_compose

all: build run

# 清理目标（如果需要）
clean:
	cd $(BACKEND_DIR) && go clean
	cd $(FRONTEND_DIR) && npm cache clean --force
