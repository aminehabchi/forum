CONTAINER = myContainer
IMAGE = forum

all: exec

build:
	@echo  "\033[1;32mBuilding Docker image... \033[0m"
	@docker image build -f Dockerfile -t $(IMAGE) .

run: build
	@clear
	@echo  "\033[1;32mRunning Docker container with Port 8080:8080... \033[0m"
	@docker container run -p 8080:8080 --detach --name $(CONTAINER) $(IMAGE)


list: run
	@echo  "\033[1;32mList All Docker images... \033[0m"
	@docker images

exec: list
	@echo  "\033[1;32mRunning Containers: \033[0m"
	@docker ps

clean:
	@docker ps -q | xargs -r docker stop
	@docker ps -a -q | xargs -r docker rm
	@docker images -q | xargs -r docker rmi
	@clear
	@echo "\033[1;32mAll Cleaned...\033[0m"