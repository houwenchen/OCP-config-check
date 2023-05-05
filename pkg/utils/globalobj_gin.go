package utils

import (
	"fmt"

	"errors"

	"github.com/gin-gonic/gin"
)

func GetConfigInfo(c *gin.Context) {
	str, err := GetInfo()
	if err != nil {
		c.JSON(400, map[string]string{
			"result": fmt.Sprintf("get config info err, %s", err),
		})
	} else {
		c.JSON(200, map[string]string{
			"result": fmt.Sprintf("get config info succeed, info: %v", str),
		})
	}
}

func SetConfig(c *gin.Context) {
	var obj GlobalObj
	err := c.ShouldBindJSON(&obj)
	if err != nil {
		err = errors.New("json format error")
		c.JSON(400, map[string]string{
			"result": fmt.Sprintf("set config err: %s", err),
		})
	} else {
		SetPara(&obj)
		c.JSON(200, map[string]string{
			"result": "set config succeed. ",
		})
	}
}
