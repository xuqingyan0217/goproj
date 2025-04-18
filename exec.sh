# 执行完毕后应当go mod tidy + go work use .

# 项目文件处于非gopath下demo_proto，所以必须指定--module为当前项目go.mod的module名，进入demo/demo_proto执行
cwgo server -I ..\..\idl\ --type RPC --module gomall/demo/demo_proto --service demo_proto --idl ..\..\idl\echo.proto

# 查看etcd容器内的数据
docker exec -it etcd etcdctl get --prefix ""

# consul的话，直接访问ui界面即可

# 生成app/frontend的home代码
cwgo server --type HTTP --idl ..\..\idl\frontend\home.proto --service frontend -module gomall/app/frontend -I ..\..\idl\

# 生成app/frontend的auth_page代码
cwgo server --type HTTP --idl ..\..\idl\frontend\auth_page.proto --service frontend -module gomall/app/frontend -I ..\..\idl\

# 生成user客户端代码
cwgo client --type RPC --service user --module gomall/rpc_gen --I ..\idl\ --idl ..\idl\user.proto

# 生成user服务端代码 --pass 是为了避免之后app里面的微服务代码对rpc_gen里面的造成影响
cwgo server --type RPC --service user --module gomall/app/user --pass "-use gomall/rpc_gen/kitex_gen" --I ..\..\idl\ --idl ..\..\idl\user.proto

# 生成product客户端代码
cwgo client --type RPC --service product --module gomall/rpc_gen --I ..\idl\ --idl ..\idl\product.proto

# 生成product服务端代码 --pass 是为了避免之后app里面的微服务代码对rpc_gen里面的造成影响
cwgo server --type RPC --service product --module gomall/app/product --pass "-use gomall/rpc_gen/kitex_gen" --I ..\..\idl\ --idl ..\..\idl\product.proto

# 生成app/frontend的product代码
cwgo server --type HTTP --idl ..\..\idl\frontend\product_page.proto --service frontend -module gomall/app/frontend -I ..\..\idl\

# 生成app/frontend的category代码
cwgo server --type HTTP --idl ..\..\idl\frontend\category_page.proto --service frontend -module gomall/app/frontend -I ..\..\idl\

# 生成cart客户端代码
cwgo client --type RPC --service cart --module gomall/rpc_gen --I ..\idl\ --idl ..\idl\cart.proto

# 生成cart服务端代码 --pass 是为了避免之后app里面的微服务代码对rpc_gen里面的造成影响
cwgo server --type RPC --service cart --module gomall/app/cart --pass "-use gomall/rpc_gen/kitex_gen" --I ..\..\idl\ --idl ..\..\idl\cart.proto

# 生成app/frontend的cart代码
cwgo server --type HTTP --idl ..\..\idl\frontend\cart_page.proto --service frontend -module gomall/app/frontend -I ..\..\idl\

# 生成payment客户端代码
cwgo client --type RPC --service payment --module gomall/rpc_gen --I ..\idl\ --idl ..\idl\payment.proto

# 生成payment服务端代码 --pass 是为了避免之后app里面的微服务代码对rpc_gen里面的造成影响
cwgo server --type RPC --service payment --module gomall/app/payment --pass "-use gomall/rpc_gen/kitex_gen" --I ..\..\idl\ --idl ..\..\idl\payment.proto

# 生成checkout客户端代码
cwgo client --type RPC --service checkout --module gomall/rpc_gen --I ..\idl\ --idl ..\idl\checkout.proto

# 生成checkout服务端代码 --pass 是为了避免之后app里面的微服务代码对rpc_gen里面的造成影响
cwgo server --type RPC --service checkout --module gomall/app/checkout --pass "-use gomall/rpc_gen/kitex_gen" --I ..\..\idl\ --idl ..\..\idl\checkout.proto

# 生成app/frontend的checkout代码
cwgo server --type HTTP --idl ..\..\idl\frontend\checkout_page.proto --service frontend -module gomall/app/frontend -I ..\..\idl\

# 生成order客户端代码
cwgo client --type RPC --service order --module gomall/rpc_gen --I ..\idl\ --idl ..\idl\order.proto

# 生成orderorder服务端代码 --pass 是为了避免之后app里面的微服务代码对rpc_gen里面的造成影响
cwgo server --type RPC --service order --module gomall/app/order --pass "-use gomall/rpc_gen/kitex_gen" --I ..\..\idl\ --idl ..\..\idl\order.proto

# 生成app/frontend的order代码
cwgo server --type HTTP --idl ..\..\idl\frontend\order_page.proto --service frontend -module gomall/app/frontend -I ..\..\idl\

# 生成email客户端代码
cwgo client --type RPC --service email --module gomall/rpc_gen --I ..\idl\ --idl ..\idl\email.proto

# 生成email服务端代码 --pass 是为了避免之后app里面的微服务代码对rpc_gen里面的造成影响
cwgo server --type RPC --service email --module gomall/app/email --pass "-use gomall/rpc_gen/kitex_gen" --I ..\..\idl\ --idl ..\..\idl\email.proto

# 生成AIEino客户端代码
cwgo client --type RPC --service AIEino --module gomall/rpc_gen --I ..\idl\ --idl ..\idl\AIEino.proto

# 生成AIEino服务端代码 --pass 是为了避免之后app里面的微服务代码对rpc_gen里面的造成影响
cwgo server --type RPC --service AIEino --module gomall/app/AIEino --pass "-use gomall/rpc_gen/kitex_gen" --I ..\..\idl\ --idl ..\..\idl\AIEino.proto

# 各个服务本身端口信息
frontend:
  addr: 0.0.0.0:8080
  Prometheus: 0.0.0.0:9998
user:
  addr: 0.0.0.0:8880
  Prometheus: 0.0.0.0:9990
product:
  addr: 0.0.0.0:8081
  Prometheus: 0.0.0.0:9991
payment:
  addr: 0.0.0.0:8083
  Prometheus: 0.0.0.0:9993
order:
  addr: 0.0.0.0:8085
  Prometheus: 0.0.0.0:9995
email:
  addr: 0.0.0.0:8086
  Prometheus: 0.0.0.0:9996
checkout:
  addr: 0.0.0.0:8084
  Prometheus: 0.0.0.0:9994
cart:
  addr: 0.0.0.0:8082
  Prometheus: 0.0.0.0:9992
aieino:
  addr: 0.0.0.0:8087
  Prometheus: 0.0.0.0:9997

kind create cluster --config=./deploy/k8s/gomall-dev-cluster.yaml

kubectl port-forward pod/xxx 8080:8080

kubectl top nodes
kind load docker-image --name gomall-dev mysql:5.7 etcd:v3.5.5 redis:7.2.5 jaegertracing/all-in-one:1.62.0 nats:2.10.19 prometheus:v2.34.0 consul:1.8.8 grafana/promtail:2.9.2 grafana/loki:2.9.2 grafana:8.3.3 alpine:3.14 nginx:1.21.4-alpine go-mall-frontend-dev:latest go-mall-product-dev:latest go-mall-email-dev:latest go-mall-checkout-dev:latest go-mall-order-dev:latest go-mall-cart-dev:latest go-mall-payment-dev:latest go-mall-user-dev:latest go-mall-aieino-dev:latest blackbox:latest alertmanager:latest go-mall-frontend-vue-dev:latest

# 生成ai的frontend代码
cwgo server --type HTTP --idl ..\..\idl\frontend\ai_eino_page.proto --service frontend -module gomall/app/frontend -I ..\..\idl\
