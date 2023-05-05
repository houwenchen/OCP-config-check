package check

import (
	"fmt"
	_ "ocp-check-config/pkg/mlog"

	"github.com/gin-gonic/gin"
)

func CheckMtu(c *gin.Context) {
	str, err := Check_Mtu()
	if err != nil {
		c.JSON(400, map[string]string{
			"result": fmt.Sprintf("check mtu failed, %s", err),
		})
	} else {
		c.JSON(200, map[string]string{
			"result": fmt.Sprintf("check mtu succeed, mtu is %s", str),
		})
	}
}

func CheckPythonEnv(c *gin.Context) {
	strs, err := Check_Python_Env()
	if err != nil {
		c.JSON(400, map[string]string{
			"result": fmt.Sprintf("check python env failed, %s", err),
		})
	} else {
		c.JSON(200, map[string]interface{}{
			"result": fmt.Sprintf("check python env succeed, package is: %s", strs),
		})
	}
}

func CheckDriver(c *gin.Context) {
	str, err := Check_Driver()
	if err != nil {
		c.JSON(400, map[string]string{
			"result": fmt.Sprintf("check ice driver failed, %s", err),
		})
	} else {
		c.JSON(200, map[string]interface{}{
			"result": fmt.Sprintf("check ice driver succeed, info: %s", str),
		})
	}
}

func CheckCapacity(c *gin.Context) {
	str, err := Check_Capacity()
	if err != nil {
		c.JSON(400, map[string]string{
			"result": fmt.Sprintf("check capacity failed, %s", err),
		})
	} else {
		c.JSON(200, map[string]string{
			"result": fmt.Sprintf("check capacity succeed, info: %s", str),
		})
	}
}

func CheckClusterNode(c *gin.Context) {
	strs, err := Check_Cluster_Node()
	if err != nil {
		c.JSON(400, map[string]string{
			"result": fmt.Sprintf("check cluster node failed, %s", err),
		})
	} else {
		c.JSON(200, map[string]string{
			"result": fmt.Sprintf("check cluster node succeed, info: %v", strs),
		})
	}
}

func CheckPTP(c *gin.Context) {
	strs, err := Check_PTP()
	if err != nil {
		c.JSON(400, map[string]string{
			"result": fmt.Sprintf("check ptp failed, %s", err),
		})
	} else {
		c.JSON(200, map[string]string{
			"result": fmt.Sprintf("check ptp succeed, info: %v", strs),
		})
	}
}
