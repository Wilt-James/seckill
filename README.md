# 电商高并发秒杀系统
## 简介
一个具备秒杀常用功能的电商系统，并根据高并发特点进行项目优化。电商秒杀系统采用分布式权限验证设计，Cookie验证的方式代替分布式Session，用Go语言开发的接口代替Redis进行数量控制防止超卖，后端采用RabbitMQ消息队列异步下单，消费消息保证Mysql数据库可用
## 运行方法
- 运行后台管理系统backend文件夹中的main.go
- 运行前台用户交互界面fronted文件夹中的main.go
- 运行分布式权限验证文件validate.go
- 运行数量控制接口文件getOne.go
- 运行异步下单数据库消费文件consumer.go
- 对于用户界面访问localhost:8082/html/htmlProduct.html
