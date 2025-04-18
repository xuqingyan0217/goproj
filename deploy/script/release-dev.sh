# 脚本列表
need_start_server_shell=(
  #rpc
  user-dev.sh
  product-dev.sh
  payment-dev.sh
  order-dev.sh
  cart-dev.sh
  email-dev.sh
  checkout-dev.sh
  aieino-dev.sh
  #http
  frontend-dev.sh
)
# 循环执行脚本
for i in ${need_start_server_shell[*]} ; do
    chmod +x $i
    ./$i
done

# 查看运行的容器，来验证是否成功
docker ps | grep go-mall

# 通过执行etcd容器的命令，来查看服务是否注册到etcd里
docker exec -it etcd etcdctl get --prefix ""

