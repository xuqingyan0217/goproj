# GoMall 微服务接口文档

## 目录

- [前端HTTP接口](#前端http接口)
  - [认证接口](#认证接口)
  - [商品接口](#商品接口)
  - [购物车接口](#购物车接口)
  - [订单接口](#订单接口)
  - [结算接口](#结算接口)
  - [AI接口](#ai接口)
- [后端RPC接口](#后端rpc接口)
  - [用户服务](#用户服务)
  - [商品服务](#商品服务)
  - [购物车服务](#购物车服务)
  - [订单服务](#订单服务)
  - [支付服务](#支付服务)
  - [结算服务](#结算服务)
  - [邮件服务](#邮件服务)
  - [AI服务](#ai服务)

## 前端HTTP接口

### 认证接口

认证相关的HTTP接口定义在 `idl/frontend/auth_page.proto` 中。

- `POST /auth/login` - 用户登录
  - 请求参数：
    - email: string (表单参数) - 用户邮箱
    - password: string (表单参数) - 用户密码
    - next: string (查询参数) - 登录后跳转地址

- `POST /auth/register` - 用户注册
  - 请求参数：
    - email: string (表单参数) - 用户邮箱
    - password: string (表单参数) - 用户密码
    - password_confirm: string (表单参数) - 确认密码

- `POST /auth/logout` - 用户登出

- `POST /auth/refresh-token` - 刷新认证令牌

### 商品接口

商品相关的HTTP接口定义在 `idl/frontend/product_page.proto` 中。

- `GET /product` - 获取商品详情
  - 请求参数：
    - id: uint32 (查询参数) - 商品ID

- `GET /search` - 搜索商品
  - 请求参数：
    - q: string (查询参数) - 搜索关键词

- `POST /product/api/create` - 创建商品
  - 请求参数：
    - name: string (表单参数) - 商品名称
    - description: string (表单参数) - 商品描述
    - price: float (表单参数) - 商品价格
    - picture: string (表单参数) - 商品图片
    - categories: []string (表单参数) - 商品分类

- `POST /product/api/update` - 更新商品
  - 请求参数：
    - id: uint32 (表单参数) - 商品ID
    - name: string (表单参数) - 商品名称
    - description: string (表单参数) - 商品描述
    - price: float (表单参数) - 商品价格
    - picture: string (表单参数) - 商品图片
    - categories: []string (表单参数) - 商品分类

- `DELETE /product/api/delete` - 删除商品
  - 请求参数：
    - id: uint32 (查询参数) - 商品ID

### 购物车接口

购物车相关的HTTP接口定义在 `idl/frontend/cart_page.proto` 中。

- `POST /cart` - 添加商品到购物车
  - 请求参数：
    - productId: uint32 (表单参数) - 商品ID
    - productNum: int32 (表单参数) - 商品数量

- `GET /cart` - 获取购物车内容

### 订单接口

订单相关的HTTP接口定义在 `idl/frontend/order_page.proto` 中。

- `GET /order` - 获取订单列表

- `POST /order/cancel` - 取消订单
  - 请求参数：
    - order_id: string - 订单ID

### 结算接口

结算相关的HTTP接口定义在 `idl/frontend/checkout_page.proto` 中。

- `GET /checkout` - 获取结算页面

- `POST /checkout/waiting` - 提交结算请求
  - 请求参数：
    - email: string (表单参数) - 用户邮箱
    - firstname: string (表单参数) - 名字
    - lastname: string (表单参数) - 姓氏
    - street: string (表单参数) - 街道地址
    - zipcode: string (表单参数) - 邮编
    - province: string (表单参数) - 省份
    - country: string (表单参数) - 国家
    - city: string (表单参数) - 城市
    - cardNum: string (表单参数) - 信用卡号
    - expirationMonth: int32 (表单参数) - 信用卡过期月份
    - expirationYear: int32 (表单参数) - 信用卡过期年份
    - cvv: int32 (表单参数) - 信用卡CVV码
    - payment: string (表单参数) - 支付方式
    - flag: uint32 (表单参数) - 标记

- `GET /checkout/result` - 获取结算结果

- `POST /checkout/repay` - 订单二次支付
  - 请求参数：
    - order_id: string - 订单ID
    - expirationMonth: int32 (表单参数) - 信用卡过期月份
    - expirationYear: int32 (表单参数) - 信用卡过期年份
    - cvv: int32 (表单参数) - 信用卡CVV码
    - payment: string (表单参数) - 支付方式
    - cardNum: string (表单参数) - 信用卡号
    - email: string - 用户邮箱

### AI接口

AI相关的HTTP接口定义在 `idl/frontend/ai_eino_page.proto` 中。

- `POST /ai/ailists` - AI查询订单列表

- `POST /ai/aiorder` - AI预下单

## 后端RPC接口

### 用户服务

用户服务的RPC接口定义在 `idl/user.proto` 中。

- `Register` - 用户注册
  - 请求参数：
    - email: string - 用户邮箱
    - password: string - 用户密码
    - password_confirm: string - 确认密码
  - 响应参数：
    - user_id: int32 - 用户ID

- `Login` - 用户登录
  - 请求参数：
    - email: string - 用户邮箱
    - password: string - 用户密码
  - 响应参数：
    - user_id: int32 - 用户ID

### 商品服务

商品服务的RPC接口定义在 `idl/product.proto` 中。

- `ListProducts` - 获取商品列表
  - 请求参数：
    - page: int32 - 页码
    - page_size: int32 - 每页数量
    - category_name: string - 分类名称
  - 响应参数：
    - products: []Product - 商品列表

- `GetProduct` - 获取商品详情
  - 请求参数：
    - id: uint32 - 商品ID
  - 响应参数：
    - product: Product - 商品信息

- `SearchProducts` - 搜索商品
  - 请求参数：
    - query: string - 搜索关键词
  - 响应参数：
    - results: []Product - 搜索结果

- `CreateProduct` - 创建商品
  - 请求参数：
    - product: Product - 商品信息
  - 响应参数：
    - id: uint32 - 商品ID

- `UpdateProduct` - 更新商品
  - 请求参数：
    - id: uint32 - 商品ID
    - product: Product - 商品信息
  - 响应参数：
    - id: uint32 - 商品ID

- `DeleteProduct` - 删除商品
  - 请求参数：
    - id: uint32 - 商品ID
  - 响应参数：
    - id: uint32 - 商品ID

- `GetAllCategory` - 获取所有分类
  - 响应参数：
    - categories: []string - 分类列表

### 购物车服务

购物车服务的RPC接口定义在 `idl/cart.proto` 中。

- `AddItem` - 添加商品到购物车
  - 请求参数：
    - user_id: uint32 - 用户ID
    - item: CartItem - 购物车商品信息
      - product_id: uint32 - 商品ID
      - quantity: uint32 - 商品数量

- `GetCart` - 获取购物车内容
  - 请求参数：
    - user_id: uint32 - 用户ID
  - 响应参数：
    - items: []CartItem - 购物车商品列表

- `EmptyCart` - 清空购物车
  - 请求参数：
    - user_id: uint32 - 用户ID

### 订单服务

订单服务的RPC接口定义在 `idl/order.proto` 中。

- `PlaceOrder` - 创建订单
  - 请求参数：
    - user_id: uint32 - 用户ID
    - user_currency: string - 用户货币
    - address: Address - 收货地址
    - email: string - 用户邮箱
    - items: []OrderItem - 订单商品列表
  - 响应参数：
    - order: OrderResult - 订单结果
      - order_id: string - 订单ID

- `ListOrder` - 获取订单列表
  - 请求参数：
    - user_id: uint32 - 用户ID
  - 响应参数：
    - orders: []Order - 订单列表

- `CancelOrder` - 取消订单
  - 请求参数：
    - order_id: string - 订单ID
    - user_id: uint32 - 用户ID
  - 响应参数：
    - success: bool - 是否成功
    - message: string - 消息

- `CancelPayment` - 取消支付
  - 请求参数：
    - transaction_id: string - 交易ID
    - order_id: string - 订单ID
    - user_id: uint32 - 用户ID
    - status: string - 状态
  - 响应参数：
    - success: bool - 是否成功
    - message: string - 消息

- `ChangeOrderStatus` - 更改订单状态
  - 请求参数：
    - order_id: string - 订单ID
    - status: string - 状态
  - 响应参数：
    - success: bool - 是否成功

### 支付服务

支付服务的RPC接口定义在 `idl/payment.proto` 中。

- `Charge` - 支付请求
  - 请求参数：
    - amount: float - 支付金额
    - credit_card: CreditCardInfo - 信用卡信息
      - credit_card_number: string - 卡号
      - credit_card_cvv: int32 - CVV码
      - credit_card_expiration_year: int32 - 过期年份
      - credit_card_expiration_month: int32 - 过期月份
    - order_id: string - 订单ID
    - user_id: uint32 - 用户ID
  - 响应参数：
    - transaction_id: string - 交易ID

### 结算服务

结算服务的RPC接口定义在 `idl/checkout.proto` 中。

- `Checkout` - 结算
  - 请求参数：
    - flag: uint32 - 标记
    - user_id: uint32 - 用户ID
    - firstname: string - 名字
    - lastname: string - 姓氏
    - email: string - 邮箱
    - address: Address - 地址信息
    - credit_card: CreditCardInfo - 信用卡信息
  - 响应参数：
    - order_id: string - 订单ID
    - transaction_id: string - 交易ID

- `PreCheckout` - 预结算
  - 请求参数：
    - user_id: uint32 - 用户ID
    - firstname: string - 名字
    - lastname: string - 姓氏
    - email: string - 邮箱
    - address: Address - 地址信息
    - product_info_list: []ProductInfo - 商品信息列表
  - 响应参数：
    - pre_order_id: string - 预订单ID
    - total_amount: float - 总金额
    - valid_until: int64 - 有效期

- `ContinueCheckout` - 继续结算
  - 请求参数：
    - order_id: string - 订单

### 邮件服务

邮件服务的RPC接口定义在 `idl/email.proto` 中。

- `SendEmail` - 发送邮件
  - 请求参数：
    - to: string - 收件人邮箱
    - subject: string - 邮件主题
    - body: string - 邮件内容
    - template_name: string - 邮件模板名称
    - template_data: map<string, string> - 模板数据
  - 响应参数：
    - success: bool - 是否发送成功
    - message: string - 返回消息

### AI服务

AI服务的RPC接口定义在 `idl/AIEino.proto` 中。

- `AIListOrder` - AI查询订单列表
  - 请求参数：
    - user_id: uint32 - 用户ID
    - query: string - 自然语言查询文本

- `AIPreOrder` - AI预下单
  - 请求参数：
    - user_id: uint32 - 用户ID
    - description: string - 商品需求描述

## 通用数据结构

### Address
地址信息结构
- street: string - 街道地址
- city: string - 城市
- province: string - 省份
- country: string - 国家
- zipcode: string - 邮编

### ProductInfo
商品信息结构
- id: uint32 - 商品ID
- name: string - 商品名称
- price: float - 商品价格
- quantity: uint32 - 商品数量
- categories: []string - 商品分类

### CreditCardInfo
信用卡信息结构
- card_number: string - 卡号
- cvv: int32 - CVV码
- expiration_year: int32 - 过期年份
- expiration_month: int32 - 过期月份

### Order
订单信息结构
- order_id: string - 订单ID
- user_id: uint32 - 用户ID
- items: []OrderItem - 订单商品列表
- total_amount: float - 订单总金额
- status: string - 订单状态
- create_time: int64 - 创建时间
- update_time: int64 - 更新时间
