VERSION=latest

# 测试环境的编译发布，编译时需要去到相应的入口文件下新建bin目录，便于存放二进制文件
build-dev:

	docker build . -f ./deploy/dockerfile/Dockerfile_alert_dev --no-cache -t alertmanager
	docker build . -f ./deploy/dockerfile/Dockerfile_blackbox_dev --no-cache -t blackbox

release-dev: build-dev


