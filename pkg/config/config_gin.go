package config

import "github.com/gin-gonic/gin"

func ConfigTenantUser(c *gin.Context) {
	DoTenantUserConfig()
}

func ConfigPyEnv(c *gin.Context) {
	DoPyConfig()
}

func ConfigAll(c *gin.Context) {
	DoConfig()
}
