#!/bin/bash

# 镜像地址
reso_addr='crpi-lofehqrjus1z8ldt.cn-beijing.personal.cr.aliyuncs.com/xqy_go/payment'
tag='latest'

container_name="go-mall-payment-dev"
# 拉取最新镜像
docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}

# 启动容器服务，其中一个端口是服务端口，一个端口是监控采集端口
docker run -p 8083:8083 -p 9993:9993 --network gomall_go_mall -v /home/logs/payment:/payment/log --name=${container_name} -d ${reso_addr}:${tag}



