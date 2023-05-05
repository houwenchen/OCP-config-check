package check

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	_ "ocp-check-config/pkg/mlog"
	"ocp-check-config/pkg/sshlib"
	"ocp-check-config/pkg/utils"
	"os"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type PTP_Profile struct {
	Interface   string `json:"interface"`
	Name        string `json:"name"`
	Phc2sysOpts string `json:"phc2sysOpts"`
	Ptp4lConf   string `json:"ptp4lConf"`
	Ptp4lOpts   string `json:"ptp4lOpts"`
}

func Check_PTP() ([]string, error) {
	//check vDU
	//先直接检查同步情况，失败了再去检查ptp同步的接口是否正确
	var ptp_pod, ptp_container string
	var strs []string
	ptp_container = "linuxptp-daemon-container"
	cmdcontext := "kubectl get pod -n openshift-ptp |grep daemon|awk '{print $1}'"

	buf, err := sshlib.ClientServe(cmdcontext)
	if err != nil {
		log.Println("sshlib.ClientServe err: ", err)
		return nil, errors.New("sshlib.ClientServe err")
	}
	ptp_pod = string(buf[:len(buf)-1])
	log.Println(ptp_pod)
	log.Println(ptp_container)

	strs = append(strs, ptp_pod)
	strs = append(strs, ptp_container)

	cmdflag := "--since=1s --limit-bytes=3000"
	cmdcontext = fmt.Sprintf("kubectl logs -f -n openshift-ptp %s -c %s %s ", ptp_pod, ptp_container, cmdflag)
	log.Println(cmdcontext)

	buf, err = sshlib.ClientServe(cmdcontext)
	if err != nil {
		log.Println("sshlib.ClientServe err: ", err)
		return nil, errors.New("sshlib.ClientServe err")
	}
	log.Println(string(buf))
	strs = append(strs, string(buf))

	info := strings.Split(string(buf), "\n")
	result := ptp4l(info[1 : len(info)-2])
	if result != nil {
		log.Println(result)
		strs = append(strs, fmt.Sprintf("%s", result))
		//trigger get ptp interface
		pp := Get_PTP_Interface()
		if pp.Interface != utils.GlobalObject.PTPInterface {
			log.Println("ptp interface is not correct")
			strs = append(strs, "ptp interface is not correct")
			//trigger change ptp interface
			Change_PTP_Interface()
			log.Println("wait 2 min, let ptp pod use new ptp config")
			time.Sleep(120 * time.Second)
			Check_PTP()
		} else {
			log.Println("ptp interface is correct, check other direction.")
			strs = append(strs, "ptp interface is correct, check other direction.")
		}
	} else {
		log.Println("ptp4l check succeed.")
		strs = append(strs, "ptp4l check succeed.")

		result = phc2sys(info[1 : len(info)-2])
		if result != nil {
			log.Println(result)
			strs = append(strs, fmt.Sprintf("%s", result))
			//trigger get ptp interface
		} else {
			log.Println("phc2sys check succeed.")
			strs = append(strs, "phc2sys check succeed.")
		}
	}
	return strs, nil
}

func ptp4l(str []string) interface{} {
	for _, v := range str {
		if !strings.Contains(v, "phc2sys") {
			if !strings.Contains(v, "offset") {
				log.Println("ptp check: OE can't get ptp4l info, please check.")
				return "ptp check: OE can't get ptp4l info, please check."
			}
		}
	}
	return nil
}

func phc2sys(str []string) interface{} {
	for _, v := range str {
		if strings.Contains(v, "phc2sys") {
			if !strings.Contains(v, "CLOCK_REALTIME") {
				log.Println("ptp check: OE can't get phc2sys info, please check.")
				return "ptp check: OE can't get phc2sys info, please check."
			}
		}
	}
	return nil
}

func Get_PTP_Interface() *PTP_Profile {
	//check OE use which interface, return interface name
	log.Println("--------start Get_PTP_Interface() scripts--------")

	pp := &PTP_Profile{}
	cmdcontext := "oc get ptpconfigs.ptp.openshift.io -n openshift-ptp -o json"
	buf, err := sshlib.ClientServe(cmdcontext)
	if err != nil {
		log.Println("sshlib.ClientServe err: ", err)
		return nil
	}
	log.Println(string(buf))
	info := gjson.Get(string(buf), "items.0.spec.profile.0")
	log.Println(info)
	err = json.Unmarshal([]byte(info.String()), pp)
	if err != nil {
		log.Println("json.Unmarshal err: ", err)
		return nil
	}
	log.Println(pp.Interface)
	log.Println("--------Get_PTP_Interface() scripts succeed--------")
	return pp
}

func Change_PTP_Interface() {
	//such as: change ens37 to ens35f0
	log.Println("--------start Change_PTP_Interface() scripts--------")

	cmdcontext := "oc get ptpconfigs.ptp.openshift.io -n openshift-ptp -o json > ./tmp/openshift-ptp.json"
	_, err := sshlib.ClientServe(cmdcontext)
	if err != nil {
		log.Println("sshlib.ClientServe error: ", err)
		return
	}

	remotefile := "/var/home/core/tmp/openshift-ptp.json"
	localfile := "openshift-ptp.json"

	//scp file to ute
	_, err = sshlib.ClientDownloadFile(remotefile, localfile)
	if err != nil {
		log.Println("sshlib.ClientDownloadFile error: ", err)
		return
	}

	//modify json file locally
	modify_ptpjson_file(localfile)

	//upload file to oe
	_, err = sshlib.ClientUploadFile(localfile, remotefile)
	if err != nil {
		log.Println("sshlib.ClientUploadFile error: ", err)
		return
	}

	log.Println("--------Change_PTP_Interface() scripts succeed--------")
	log.Println("--------start commit new ptp interface--------")

	cmdcontext = "oc apply -f /var/home/core/tmp/openshift-ptp.json"
	buf, err := sshlib.ClientServe(cmdcontext)
	if err != nil {
		log.Println("sshlib.ClientServe error: ", err)
		return
	}
	log.Println(string(buf))

	log.Println("--------commit new ptp interface succeed--------")
}

//modify ens37 to ens35f0
/*
    "interface": "ens37",
    "name": "ptp-profile",
    "phc2sysOpts": "-n 24 -i ens37 -R 1 -w",
    "ptp4lConf": "...",
    "ptp4lOpts": "-2 -s"

	to

    "interface": "ens35f0",
    "name": "ptp-profile",
    "phc2sysOpts": "-n 24 -i ens35f0 -R 1 -w",
    "ptp4lConf": "...",
    "ptp4lOpts": "-2 -s"
*/
func modify_ptpjson_file(path string) {
	buf, err := os.ReadFile(path)
	if err != nil {
		log.Println("os.ReadFile error: ", err)
		return
	}

	str, err := sjson.Set(string(buf), "items.0.spec.profile.0.interface", utils.GlobalObject.PTPInterface)
	if err != nil {
		log.Println("sjson.Set error: ", err)
		return
	}

	info := fmt.Sprintf("-n 24 -i %s -R 1 -w", utils.GlobalObject.PTPInterface)
	str, err = sjson.Set(str, "items.0.spec.profile.0.phc2sysOpts", info)
	if err != nil {
		log.Println("sjson.Set error: ", err)
		return
	}

	err = os.WriteFile(path, []byte(str), 0666)
	if err != nil {
		log.Println("os.WriteFile error: ", err)
		return
	}
}
