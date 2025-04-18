# 脚本列表
need_start_server_shell=(
  # rpc
  user-rpc-test.sh
  social-rpc-test.sh
  im-rpc-test.sh
  im-ws-test.sh
  # api
  user-api-test.sh
  social-api-test.sh
  im-api-test.sh
  # task
  task-mq-test.sh
)
# 循环执行脚本
for i in ${need_start_server_shell[*]} ; do
    chmod +x $i
    ./$i
done

# 查看运行的容器，来验证是否成功
docker ps | grep easy-chat

# 通过执行etcd容器的命令，来查看服务是否注册到etcd里
docker exec -it etcd etcdctl get --prefix ""


