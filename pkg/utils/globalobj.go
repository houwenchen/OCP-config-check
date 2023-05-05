package utils

import (
	"encoding/json"
	"errors"
	"log"
	_ "ocp-check-config/pkg/mlog"
	"os"
)

type GlobalObj struct {
	/*
		OE info
	*/
	SshIP         string `json:"sship"`
	User          string `json:"user"`
	Passwd        string `json:"password"`
	TenantUserNum int    `json:"tenantusernumber"`
	PTPInterface  string `json:"ptpinterface"`
	NeType        string `json:"netype"` //options: vDU/vCU

	/*
		download infra-oam scripts url
	*/
	ScriptsVersion string `json:"scriptsversion"`

	/*
		config file path
	*/
	ConfFilePath string `json:"conffilepath"`
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		SshIP:         "10.69.72.110",
		User:          "core",
		Passwd:        "system123",
		TenantUserNum: 1,
		PTPInterface:  "ens35f0",
		NeType:        "vDU",
		ConfFilePath:  "/home/HZtcuser/tmp/go/src/OCP_with_Command/conf/config.json",
	}
	GlobalObject.LoadConf()
}

func (g *GlobalObj) LoadConf() {
	data, err := os.ReadFile(g.ConfFilePath)
	if err != nil {
		log.Println("os.ReadFile err: ", err)
		return
	}
	err = json.Unmarshal(data, g)
	if err != nil {
		log.Println("json.Unmarshal err: ", err)
		return
	}
}

func (g *GlobalObj) CompletePara() {
	if g.SshIP == "" {
		g.SshIP = "0.0.0.0"
	}
	if g.User == "" {
		g.User = "core"
	}
	if g.Passwd == "" {
		g.Passwd = "system123"
	}
	if g.TenantUserNum == 0 {
		g.TenantUserNum = 1
	}
	if g.PTPInterface == "" {
		g.PTPInterface = "ens43"
	}
	if g.NeType == "" {
		g.NeType = "vDU"
	}
	if g.ScriptsVersion == "" {
		g.ScriptsVersion = "2.27.0"
	}
	if g.ConfFilePath == "" {
		g.ConfFilePath = "./conf/config.json"
	}
}

func (g *GlobalObj) ReloadConf(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("os.ReadFile err: ", err)
		return
	}
	err = json.Unmarshal(data, g)
	if err != nil {
		log.Println("json.Unmarshal err: ", err)
		return
	}
}

func GetInfo() (string, error) {
	byteData, err := json.Marshal(GlobalObject)
	if err != nil {
		log.Println("json.Marshal err: ", err)
		return "", errors.New("json.Marshal err")
	}
	log.Println(string(byteData))
	return string(byteData), nil
}

// 设置参数：IP，user, passwd, ptp
func SetPara(obj *GlobalObj) {
	GlobalObject = obj
}
