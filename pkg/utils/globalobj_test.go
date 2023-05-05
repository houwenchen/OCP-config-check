package utils

import "testing"

func TestLoadConf(t *testing.T) {
	g := &GlobalObj{ConfFilePath: "../../conf/config.json"}
	g.LoadConf()
	if g.SshIP != "10.69.72.110" || g.User != "core" || g.Passwd != "system123" || g.NeType != "vDU" {
		panic("test failed.")
	}
}
