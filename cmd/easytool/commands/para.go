package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"ocp-check-config/pkg"
	_ "ocp-check-config/pkg/mlog"
	"ocp-check-config/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	sshIp            string
	user             string
	password         string
	tenantusernumber int
	netype           string
	ptpinterface     string
	scriptsversion   string
)

var ParaCmd = &cobra.Command{
	Use:   "para",
	Short: "operation about para. ",
	Long:  "operation about para. ",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("operation about para.")
	},
}

var ParaGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get config info. ",
	Long:  "get config info. ",
	Run: func(cmd *cobra.Command, args []string) {
		err := BaseRequestForPara("getconfiginfo")
		if err != nil {
			log.Println(err)
		}
	},
}

var ParaSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set config info. ",
	Long:  "set config info. ",
	Run: func(cmd *cobra.Command, args []string) {
		err := SendSetParaRequest(sshIp, user, password, tenantusernumber, netype, ptpinterface, scriptsversion)
		if err != nil {
			log.Println(err)
		}
	},
}

func BaseRequestForPara(postfix string) error {
	cli := &http.Client{}
	serverUrl := "http://" + "localhost:9090" + pkg.BaseURL + postfix
	req, err := http.NewRequest("GET", serverUrl, nil)
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

func SendSetParaRequest(ip string, user string, passwd string, tunum int, netype string, ptpi string, sv string) error {
	obj := &utils.GlobalObj{
		SshIP:          ip,
		User:           user,
		Passwd:         passwd,
		TenantUserNum:  tunum,
		NeType:         netype,
		PTPInterface:   ptpi,
		ScriptsVersion: sv,
	}
	byteData, err := json.Marshal(obj)
	if err != nil {
		log.Println("json.Marshal err: ", err)
		return err
	}

	cli := &http.Client{}
	serverUrl := "http://" + "localhost:9090" + pkg.BaseURL + "setpara"
	req, err := http.NewRequest("POST", serverUrl, bytes.NewReader(byteData))
	if err != nil {
		log.Println("http.NewRequest err: ", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := cli.Do(req)
	if err != nil {
		log.Println("cli.Do err: ", err)
		return err
	}
	defer resp.Body.Close()

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
	RootCmd.AddCommand(ParaCmd)
	ParaCmd.AddCommand(ParaGetCmd)
	ParaCmd.AddCommand(ParaSetCmd)
	ParaSetCmd.Flags().StringVarP(&sshIp, "sship", "", "10.0.0.1", "remote ssh ip. ")
	ParaSetCmd.Flags().StringVarP(&user, "user", "u", "core", "remote login user. ")
	ParaSetCmd.Flags().StringVarP(&password, "password", "p", "system123", "remote login password. ")
	ParaSetCmd.Flags().IntVarP(&tenantusernumber, "tenantusernumber", "n", 1, "tenant user number of remote host. ")
	ParaSetCmd.Flags().StringVarP(&netype, "netype", "t", "vDU", "netype of remote host. ")
	ParaSetCmd.Flags().StringVarP(&ptpinterface, "ptpinterface", "i", "ens35f0", "ptp interface of remote host. ")
	ParaSetCmd.Flags().StringVarP(&scriptsversion, "scriptsversion", "s", "2.27.0", "infra-oam tool version. ")
}
