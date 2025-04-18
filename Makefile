# 所有的命令都要在根目录下执行，命令最终落脚到每个脚本上时，其路径都是以根目录为起始的
# release-dev一个是入口启动名称，一个是mk文件里面的名称，二者不一样
frontend-dev:
	@make -f deploy/mk/frontend-dev.mk release-dev

product-dev:
	@make -f deploy/mk/product-dev.mk release-dev

user-dev:
	@make -f deploy/mk/user-dev.mk release-dev

payment-dev:
	@make -f deploy/mk/payment-dev.mk release-dev

cart-dev:
	@make -f deploy/mk/cart-dev.mk release-dev

order-dev:
	@make -f deploy/mk/order-dev.mk release-dev

checkout-dev:
	@make -f deploy/mk/checkout-dev.mk release-dev

email-dev:
	@make -f deploy/mk/email-dev.mk release-dev

aieino-dev:
	@make -f deploy/mk/aieino-dev.mk release-dev

# 统一打包后端镜像入口
release-dev: frontend-dev product-dev user-dev payment-dev cart-dev order-dev checkout-dev email-dev aieino-dev


# 统一启动入口,ubuntu需要加bash,centos不需要
# 仅针对docker部署
install-server:
	cd ./deploy/script && chmod +x release-dev.sh && bash ./release-dev.sh
# prometheus-exports 组件镜像
release-prometheus-exports:
	@make -f deploy/mk/prometheus-exports.mk release-dev

install-frontend-dev:
	cd ./deploy/script && chmod +x frontend-dev.sh && bash ./frontend-dev.sh

install-product-dev:
	cd ./deploy/script && chmod +x product-dev.sh && bash ./product-dev.sh

install-user-devdev:
	cd ./deploy/script && chmod +x user-dev.sh && bash ./user-dev.sh

install-cart-dev:
	cd ./deploy/script && chmod +x cart-dev.sh && bash ./cart-dev.sh

install-order-dev:
	cd ./deploy/script && chmod +x order-dev.sh && bash ./order-dev.sh

install-payment-dev:
	cd ./deploy/script && chmod +x payment-dev.sh && bash ./payment-dev.sh

install-checkout-dev:
	cd ./deploy/script && chmod +x checkout-dev.sh && bash ./checkout-dev.sh

install-email-dev:
	cd ./deploy/script && chmod +x email-dev.sh && bash ./email-dev.sh

install-aieino-dev:
	cd ./deploy/script && chmod +x aieino-dev.sh && bash ./aieino-dev.sh


