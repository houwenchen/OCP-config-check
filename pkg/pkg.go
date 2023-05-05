package pkg

import (
	"ocp-check-config/pkg/check"
	"ocp-check-config/pkg/config"
	"ocp-check-config/pkg/utils"

	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
}

var BaseURL = "/api/easytool/v1/"

func NewApp() *App {
	app := &App{
		engine: gin.Default(),
	}
	routeGroup := app.engine.Group(BaseURL)

	//checkall 使用中间件的方法
	routeGroup.GET("checkall", check.CheckMtu, check.CheckPythonEnv, check.CheckDriver, check.CheckCapacity, check.CheckClusterNode, check.CheckPTP)
	routeGroup.GET("checkmtu", check.CheckMtu)
	routeGroup.GET("checkpyenv", check.CheckPythonEnv)
	routeGroup.GET("checkdriver", check.CheckDriver)
	routeGroup.GET("checkcapacity", check.CheckCapacity)
	routeGroup.GET("checkclusternode", check.CheckClusterNode)
	routeGroup.GET("checkptp", check.CheckPTP)

	routeGroup.GET("getconfiginfo", utils.GetConfigInfo)
	routeGroup.POST("setpara", utils.SetConfig)

	routeGroup.PUT("configall", config.ConfigAll)
	routeGroup.PUT("configtenantuser", config.ConfigTenantUser)
	routeGroup.PUT("configpyenv", config.ConfigPyEnv)

	return app
}

func (a *App) Run() error {
	err := a.engine.Run(":9090")
	return err
}
