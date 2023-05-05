package check

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	_ "ocp-check-config/pkg/mlog"
	"ocp-check-config/pkg/sshlib"
	"strings"

	"github.com/tidwall/gjson"
)

func Check_All() {
	Check_Python_Env()
	Check_Mtu()
	Check_Driver()
	Check_Capacity()
	Check_PTP()
	Check_Cluster_Node()
}

func Check_Mtu() (string, error) {
	//mtu should be 1500
	log.Println("--------start check_mtu() scripts--------")

	cmdcontext := "netstat -i|grep br-ex|awk '{print $2}'"
	buf, err := sshlib.ClientServe(cmdcontext)
	if err != nil {
		log.Println("sshlib.ClientServe err: ", err)
		return "", errors.New("sshlib.ClientServe err")
	}
	mtu := string(buf[:len(buf)-1])
	log.Printf("--------mtu info is %s--------\n", mtu)
	// result := strings.Contains(mtu, "1500")
	// if !result {
	// 	fmt.Println("mtu is not equal 1500")
	// 	return
	// }
	result := strings.Compare(mtu, "1500")
	if result != 0 {
		log.Println("mtu is not equal 1500")
		return mtu, errors.New("mtu is not equal 1500. ")
	}

	log.Println("--------check_mtu() scripts succeed--------")
	return mtu, nil
}

func Check_Driver() (string, error) {
	//driver should be ice: 1.3.2
	log.Println("--------start check_driver() scripts--------")

	net_inf := Get_Net_Inf()
	cmdcontext := "ethtool -i " + net_inf[0]
	buf, err := sshlib.ClientServe(cmdcontext)
	if err != nil {
		log.Println("sshlib.ClientServe err: ", err)
		return "", err
	}
	if !strings.Contains(string(buf), "driver: ice") {
		log.Println("driver is not ice")
		return "", errors.New("driver is not ice")
	}
	if !strings.Contains(string(buf), "version: 1.3.2") {
		log.Println("driver version is not correct")
		return "", errors.New("driver version is not correct")
	}

	log.Println("--------check_driver() scripts succeed--------")
	return string(buf), nil
}

func Get_Net_Inf() []string {
	cmdcontext := "netstat -i|grep ens|awk '{print $1}'"
	buf, err := sshlib.ClientServe(cmdcontext)
	if err != nil {
		log.Println("sshlib.ClientServe err: ", err)
		return nil
	}
	inf := strings.Fields(string(buf))
	return inf
}

func Check_Capacity() (string, error) {
	log.Println("--------start check_capacity() scripts--------")

	netdevice_num := "30"
	vfio_num := "20"

	net_inf := Get_Net_Inf()
	netdevice_info := make([]string, 0)
	vfio_info := make([]string, 0)
	for _, v := range net_inf {
		netdevice_info = append(netdevice_info, fmt.Sprintf("openshift.io/sriov_netdevice_%s", v))
		vfio_info = append(vfio_info, fmt.Sprintf("openshift.io/sriov_vfio_%s", v))
	}

	cmdcontext := "kubectl get node  -o json"
	buf, err := sshlib.ClientServe(cmdcontext)
	if err != nil {
		log.Println("sshlib.ClientServe err: ", err)
		return "", errors.New("sshlib.ClientServe err")
	}
	// fmt.Println(string(buf))

	info := gjson.Get(string(buf), "items.0.status.capacity")
	//反序列化信息到map结构中
	anyMap := make(map[string]interface{}, 0)
	if err := json.Unmarshal([]byte(info.String()), &anyMap); err != nil {
		log.Println("json.Unmarshal error: ", err)
		return "", errors.New("json.Unmarshal error")
	}

	for _, v := range netdevice_info {
		if anyMap[v] != netdevice_num {
			log.Printf("%s is not in correct num: %s", v, netdevice_num)
			return "", errors.New("netdevice have incorrect num")
		}
	}

	for _, v := range vfio_info {
		if anyMap[v] != vfio_num {
			log.Printf("%s is not in correct num: %s", v, netdevice_num)
			return "", errors.New("vfio have incorrect num")
		}
	}
	log.Println("--------check_capacity() succeed--------")
	return "all netdevice = 30, all vfio = 20", nil
}

func Check_Cluster_Node() (interface{}, error) {
	//check cluster node status
	//return: error or cluster_node_info map
	/*
		NAME: hztt-10-69-72-110
		STATUS: Ready
		ROLES: master,worker
		AGE: 176d
		VERSION: v1.23.5+9ce5071
	*/
	log.Println("--------start Check_Cluster_Node() scripts--------")

	cmdcontext := "kubectl get node"
	buf, err := sshlib.ClientServe(cmdcontext)
	if err != nil {
		log.Println("sshlib.ClientServe err: ", err)
		return nil, errors.New("sshlib.ClientServe err")
	}
	inf := strings.Fields(string(buf))
	lenth := len(inf)
	anyMap := make(map[string]string, lenth/2)
	for i := 0; i < lenth/2; i++ {
		key := inf[i]
		value := inf[i+lenth/2]
		anyMap[key] = value
		log.Printf("%s: %s\n", key, value)
	}

	log.Println("--------Check_Cluster_Node() scripts succeed--------")

	return anyMap, nil
}

func Check_Python_Env() ([]string, error) {
	//check python env, such as: requests, pyyaml==5.4.1
	log.Println("--------start Check_Python_Env() scripts--------")

	cmdcontext := "python -m pip list"
	buf, err := sshlib.ClientServe(cmdcontext)
	if err != nil {
		log.Println("sshlib.ClientServe err: ", err)
		return nil, err
	}
	if !strings.Contains(string(buf), "requests") {
		log.Println("python env lose requests package.")
		return nil, errors.New("python env lose requests package")
	}

	if !strings.Contains(string(buf), "PyYAML (5.4.1)") {
		log.Println("python env lose PyYAML package.")
		return nil, errors.New("python env lose PyYAML package")
	}

	inf := strings.Split(string(buf), "\n")
	newinf := make([]string, 0)
	for _, value := range inf {
		if !strings.Contains(value, "DEPRECATION") {
			newinf = append(newinf, value)
		}
	}
	log.Println(newinf)
	log.Println("--------Check_Python_Env() scripts succeed--------")
	return newinf, nil
}
