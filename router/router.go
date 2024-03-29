package router

import (
	"fpga-bms-server/config"
	"fpga-bms-server/controllers"
	"fpga-bms-server/database"
	"fpga-bms-server/middleware"
	"fpga-bms-server/models"
	"github.com/betacraft/yaag/yaag"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/middleware/logger"
)

var (
	Router = newApp()
)

/**
 * 初始化 iris app
 * @method NewApp
 * @return  {[type]}  api      *iris.Application  [iris app]
 */
func newApp() (api *iris.Application) {
	api = iris.New()
	api.Use(logger.New())

	api.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.JSON(controllers.ApiResource(false, nil, "404 Not Found"))
	})
	api.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.WriteString("Oups something went wrong, try again")
	})

	//同步模型数据表
	//如果模型表这里没有添加模型，单元测试会报错数据表不存在。
	//因为单元测试结束，会删除数据表
	database.DB.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
	)

	iris.RegisterOnInterrupt(func() {
		database.DB.Close()
	})

	//"github.com/iris-contrib/middleware/cors"
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})

	appName := config.Conf.Get("app.name").(string)
	appDoc := config.Conf.Get("app.doc").(string)
	appUrl := config.Conf.Get("app.url").(string)
	//api 文档配置
	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
		On:       true,
		DocTitle: appName,
		DocPath:  appDoc + "/index.html", //设置绝对路径
		BaseUrls: map[string]string{
			"Production": appUrl,
			"Staging":    "",
		},
	})

	//初始化系统 账号 权限 角色
	models.CreateSystemData()

	v1 := api.Party("/v1", crs).AllowMethods(iris.MethodOptions)
	{
		v1.Use(middleware.NewYaag()) // <- IMPORTANT, register the middleware.
		v1.Post("/admin/login", controllers.UserLogin)
		v1.PartyFunc("/admin", func(admin router.Party) {
			admin.Use(middleware.JwtHandler().Serve, middleware.AuthToken)
			admin.Get("/logout", controllers.UserLogout)

			admin.PartyFunc("/users", func(users router.Party) {
				users.Get("/", controllers.GetAllUsers)
				users.Get("/{id:uint}", controllers.GetUser)
				users.Post("/", controllers.CreateUser)
				users.Put("/{id:uint}", controllers.UpdateUser)
				users.Delete("/{id:uint}", controllers.DeleteUser)
				users.Get("/profile", controllers.GetProfile)
			})
			
			admin.PartyFunc("/roles", func(roles router.Party) {
				roles.Get("/", controllers.GetAllRoles)
				roles.Get("/{id:uint}", controllers.GetRole)
				roles.Post("/", controllers.CreateRole)
				roles.Put("/{id:uint}", controllers.UpdateRole)
				roles.Delete("/{id:uint}", controllers.DeleteRole)
			})
			admin.PartyFunc("/permissions", func(permissions router.Party) {
				permissions.Get("/", controllers.GetAllPermissions)
				permissions.Get("/{id:uint}", controllers.GetPermission)
				permissions.Post("/", controllers.CreatePermission)
				permissions.Put("/{id:uint}", controllers.UpdatePermission)
				permissions.Delete("/{id:uint}", controllers.DeletePermission)
			})
			
		})
	}

	return
}
