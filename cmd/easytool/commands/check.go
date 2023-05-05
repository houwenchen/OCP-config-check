package cmd

import (
	"io"
	"log"
	"net/http"
	"ocp-check-config/pkg"
	_ "ocp-check-config/pkg/mlog"

	"github.com/spf13/cobra"
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "check oe's env. ",
	Long:  "check oe's env. ",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("check oe's env. ")
	},
}

var CheckAllCmd = &cobra.Command{
	Use:   "all",
	Short: "check all info. ",
	Long:  "check all info. ",
	Run: func(cmd *cobra.Command, args []string) {
		err := BaseRequestForAll("checkall")
		if err != nil {
			log.Println(err)
		}
	},
}

var CheckMtuCmd = &cobra.Command{
	Use:   "mtu",
	Short: "check mtu ?= 1500, if != 1500, will meet ssh or sctp failed. ",
	Long:  "check mtu ?= 1500, if != 1500, will meet ssh or sctp failed. ",
	Run: func(cmd *cobra.Command, args []string) {
		err := BaseRequest("checkmtu")
		if err != nil {
			log.Println(err)
		}
	},
}

var CheckPyEnvCmd = &cobra.Command{
	Use:   "pyenv",
	Short: "check whether python env contains requests package and PyYAML package. ",
	Long:  "check whether python env contains requests package and PyYAML package. ",
	Run: func(cmd *cobra.Command, args []string) {
		err := BaseRequest("checkpyenv")
		if err != nil {
			log.Println(err)
		}
	},
}

var CheckDriverCmd = &cobra.Command{
	Use:   "driver",
	Short: "check whether ice driver version is 1.3.2. ",
	Long:  "check whether ice driver version is 1.3.2. ",
	Run: func(cmd *cobra.Command, args []string) {
		err := BaseRequest("checkdriver")
		if err != nil {
			log.Println(err)
		}
	},
}

var CheckCapacityCmd = &cobra.Command{
	Use:   "capacity",
	Short: "check netdevice num and vfio num of ervery network interface. ",
	Long:  "check netdevice num and vfio num of ervery network interface. ",
	Run: func(cmd *cobra.Command, args []string) {
		err := BaseRequest("checkcapacity")
		if err != nil {
			log.Println(err)
		}
	},
}

var CheckClusterNodeCmd = &cobra.Command{
	Use:   "clusternode",
	Short: "get cluster node info. ",
	Long:  "get cluster node info. ",
	Run: func(cmd *cobra.Command, args []string) {
		err := BaseRequest("checkclusternode")
		if err != nil {
			log.Println(err)
		}
	},
}

var CheckPTPCmd = &cobra.Command{
	Use:   "ptp",
	Short: "check cluster's ptp status. ",
	Long:  "check cluster's ptp status. ",
	Run: func(cmd *cobra.Command, args []string) {
		err := BaseRequest("checkptp")
		if err != nil {
			log.Println(err)
		}
	},
}

func BaseRequest(postfix string) error {
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

func BaseRequestForAll(postfix string) error {
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
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println("io.ReadAll err: ", err)
		return err
	}
	// log.Println(string(respbyte))
	return nil
}

func init() {
	RootCmd.AddCommand(CheckCmd)
	CheckCmd.AddCommand(CheckAllCmd)
	CheckCmd.AddCommand(CheckMtuCmd)
	CheckCmd.AddCommand(CheckPyEnvCmd)
	CheckCmd.AddCommand(CheckDriverCmd)
	CheckCmd.AddCommand(CheckCapacityCmd)
	CheckCmd.AddCommand(CheckClusterNodeCmd)
	CheckCmd.AddCommand(CheckPTPCmd)
}
