package cmd

import (
	"io"
	"log"
	"net/http"
	"ocp-check-config/pkg"
	_ "ocp-check-config/pkg/mlog"

	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "do config operation. ",
	Long:  "do config operation. ",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("do config operation.")
	},
}

var ConfigAllCmd = &cobra.Command{
	Use:   "all",
	Short: "set all env. ",
	Long:  "set all env. ",
	Run: func(cmd *cobra.Command, args []string) {
		err := ConfigBaseRequest("configall")
		if err != nil {
			log.Println(err)
		}
	},
}

var ConfigPyCmd = &cobra.Command{
	Use:   "py",
	Short: "set python env. ",
	Long:  "set python env. ",
	Run: func(cmd *cobra.Command, args []string) {
		err := ConfigBaseRequest("configpyenv")
		if err != nil {
			log.Println(err)
		}
	},
}

var ConfigTenantUserCmd = &cobra.Command{
	Use:   "tenantuser",
	Short: "set tenantuser env. ",
	Long:  "set tenantuser env. ",
	Run: func(cmd *cobra.Command, args []string) {
		err := ConfigBaseRequest("configtenantuser")
		if err != nil {
			log.Println(err)
		}
	},
}

func ConfigBaseRequest(postfix string) error {
	cli := &http.Client{}
	serverUrl := "http://" + "localhost:9090" + pkg.BaseURL + postfix
	req, err := http.NewRequest("PUT", serverUrl, nil)
	if err != nil {
		log.Println("http.NewRequest err: ", err)
		return err
	}
	resp, err := cli.Do(req)
	if err != nil {
		log.Println("client.Do err: ", err)
		return err
	}
	log.Println(resp.Status)
	respbyte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("io.ReadAll err: ", err)
		return err
	}
	log.Println(string(respbyte))
	return nil
}

func init() {
	RootCmd.AddCommand(ConfigCmd)
	ConfigCmd.AddCommand(ConfigAllCmd)
	ConfigCmd.AddCommand(ConfigPyCmd)
	ConfigCmd.AddCommand(ConfigTenantUserCmd)
}
