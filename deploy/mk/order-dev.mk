VERSION=latest
SERVER_NAME=order
SERVER_TYPE=dev
# 测试环境配置
# docker的镜像发布地址，这里我们改为我们的私有仓库
DOCKER_REPO_TEST=crpi-lofehqrjus1z8ldt.cn-beijing.personal.cr.aliyuncs.com/xqy_go/${SERVER_NAME}
# 测试版本
VERSION_TEST=$(VERSION)
# 编译的程序名称
APP_NAME_TEST=go-mall-${SERVER_NAME}-${SERVER_TYPE}
# 测试下的编译文件
DOCKER_FILE_TEST=./deploy/dockerfile/Dockerfile_${SERVER_NAME}_${SERVER_TYPE}
# 测试环境的编译发布，编译时需要去到相应的入口文件下新建bin目录，便于存放二进制文件
build-dev:

	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${SERVER_NAME}-${SERVER_TYPE} ./app/${SERVER_NAME}/
	docker build . -f ${DOCKER_FILE_TEST} --no-cache -t ${APP_NAME_TEST}

# 镜像的测试标签
tag-dev:

	@echo 'create tag ${VERSION_TEST}'
	docker tag ${APP_NAME_TEST} ${DOCKER_REPO_TEST}:${VERSION_TEST}

publish-dev:

	@echo 'publish ${VERSION_TEST} to ${DOCKER_REPO_TEST}'
	docker push $(DOCKER_REPO_TEST):${VERSION_TEST}

release-dev: build-dev tag-dev publish-dev


