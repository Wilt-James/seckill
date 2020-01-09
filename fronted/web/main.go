package main

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/opentracing/opentracing-go/log"
	"seckill/common"
	"seckill/fronted/middleware"
	"seckill/fronted/web/controllers"
	"seckill/rabbitmq"
	"seckill/repositories"
	"seckill/services"
)

func main() {
	//1. 创建iris实例
	app := iris.New()
	//2. 设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3. 注册模板
	tpl := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tpl)
	//4. 设置模板目标
	app.StaticWeb("/public", "./fronted/web/public")
	//访问生成好的html静态文件，nginx配合实现
	app.StaticWeb("/html", "./fronted/web/htmlProductShow")
	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	// 连接至数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	/*sess := sessions.New(sessions.Config{
		Cookie:   "AdminCookie",
		Expires:  600 * time.Minute,
	})*/
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//5. 注册控制器
	userRepository := repositories.NewUserManager("users", db)
	userService := services.NewUserService(userRepository)
	userParty := app.Party("/user")
	user := mvc.New(userParty)
	user.Register(ctx, userService)
	user.Handle(new(controllers.UserController))

	rabbitmq := rabbitmq.NewRabbitMQSimple("seckill_product")

	// 注册product控制器
	productRepository := repositories.NewProductManager("products", db)
	productService := services.NewProductService(productRepository)
	orderRepository := repositories.NewOrderManager("orders", db)
	orderService := services.NewOrderService(orderRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	productParty.Use(middleware.AuthConProduct)
	product.Register(productService, orderService, rabbitmq)
	product.Handle(new(controllers.ProductController))



	//6. 启动服务
	app.Run(
		iris.Addr("localhost:8082"),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)


}
