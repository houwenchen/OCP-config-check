package config

import (
	"fmt"
	"log"
	_ "ocp-check-config/pkg/mlog"
	"ocp-check-config/pkg/utils"
	"os"

	"gopkg.in/ini.v1"
)

/*
需求：

	使用ini库实现对hosts文件的直接处理
*/
func modify_hosts(path, newpath string) {
	//path 是原文件路径， newpath是修改后的文件保存路径
	cfg, err := ini.Load(path)
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	cfg.Section("ocp").NewBooleanKey(utils.GlobalObject.SshIP)

	cfg.Section("ocp:vars").Key("ansible_user").SetValue(utils.GlobalObject.User)
	log.Println("after modify, ansible_user: ", cfg.Section("ocp:vars").Key("ansible_user").String())

	cfg.Section("ocp:vars").Key("ansible_password").SetValue(utils.GlobalObject.Passwd)
	log.Println("after modify, ansible_password: ", cfg.Section("ocp:vars").Key("ansible_password").String())

	err = cfg.SaveTo(newpath)
	if err != nil {
		log.Println("cfg.SaveTo error: ", err)
		panic(fmt.Sprintf("cfg.SaveTo error: %v", err))
	}
	log.Printf("cfg.SaveTo %s succeed.", newpath)
}
