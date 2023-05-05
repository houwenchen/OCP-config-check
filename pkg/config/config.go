package config

import (
	"fmt"
	"log"
	_ "ocp-check-config/pkg/mlog"
	"ocp-check-config/pkg/sshlib"
	"ocp-check-config/pkg/utils"
	"os"
	"os/exec"
)

func DoConfig() {
	download_scripts()
	tar_scripts()
	modify_config()
	create_tenant_user()
	python_env()
	python_env_other()
}

func DoPyConfig() {
	download_scripts()
	tar_scripts()
	modify_config()
	python_env()
	python_env_other()
}

func DoTenantUserConfig() {
	download_scripts()
	tar_scripts()
	modify_config()
	create_tenant_user()
}

func PreConfigEnv() {
	download_scripts()
	tar_scripts()
	modify_config()
}

func command_run_check(cmdcontext string) {
	cmd := exec.Command("bash", "-c", cmdcontext)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Println("command exec failed, err :", err)
	}
}

func download_scripts() {
	//download infra-oam scripts
	log.Println("--------start download infra-oam scripts--------")
	cmd := "rm -rf infra-om*"
	command_run_check(cmd)
	//https://gitlabe2.ext.net.nokia.com/cloud-infra/infra-om/-/archive/2.27.0/infra-om-2.27.0.tar.gz
	url := "https://gitlabe2.ext.net.nokia.com/cloud-infra/infra-om/-/archive/" + utils.GlobalObject.ScriptsVersion + "/infra-om-" + utils.GlobalObject.ScriptsVersion + ".tar.gz"
	cmdcontext := "wget " + url
	command_run_check(cmdcontext)
}

func tar_scripts() {
	//check infra-oam scripts exists and tar files
	log.Println("--------start check infra-oam scripts files and tar files--------")
	fileName := "infra-om-" + utils.GlobalObject.ScriptsVersion + ".tar.gz"
	cmdcontext := "tar -xvzf " + fileName
	command_run_check(cmdcontext)
}

func modify_config() {
	//modify infra-om-2.26.0/hosts fils, set hosts, user and password
	//-----CAHNGELOG: use config_hosts scripts
	// log.Println("--------start modify infra-om-2.26.0/hosts fils, set hosts, user and password--------")
	// set_ip_info := fmt.Sprintf("\"2i %s\"", utils.GlobalObject.SshIP)
	// set_user_info := fmt.Sprintf("\"s/ansible_user=core/ansible_user=%s/g\"", utils.GlobalObject.User)
	// set_password_info := fmt.Sprintf("\"s/ansible_password=/ansible_password=%s/g\"", utils.GlobalObject.Passwd)
	// cmdcontext := "sed -i " + set_ip_info + " infra-om-2.26.0/hosts"
	// command_run_check(cmdcontext)
	// cmdcontext = "sed -i " + set_user_info + " infra-om-2.26.0/hosts"
	// command_run_check(cmdcontext)
	// cmdcontext = "sed -i " + set_password_info + " infra-om-2.26.0/hosts"
	// command_run_check(cmdcontext)
	// cmdcontext = "cat infra-om-2.26.0/hosts"
	// command_run_check(cmdcontext)
	fileName := "infra-om-" + utils.GlobalObject.ScriptsVersion
	path := fileName + "/hosts"
	newpath := fileName + "/hosts"
	modify_hosts(path, newpath)
}

func create_tenant_user() {
	v1 := utils.GlobalObject.ScriptsVersion
	v2 := utils.GlobalObject.TenantUserNum
	log.Println("--------start create tenant user via infra-oam scripts--------")
	cmdcontext := fmt.Sprintf("ansible-playbook -i infra-om-%s/hosts infra-om-%s/configure-tenant.yml -e '{first_user_index: 1, last_user_index: %v, configure_action: create,  test_env: 5G, root_flag: enable}'", v1, v1, v2)
	command_run_check(cmdcontext)
	log.Println(cmdcontext)
	cmdcontext = fmt.Sprintf("ansible-playbook -i infra-om-%s/hosts infra-om-%s/configure-sudo-commands.yml -e '{first_user_index: 1, last_user_index: %v, test_env: 5G}'", v1, v1, v2)
	command_run_check(cmdcontext)
	log.Println(cmdcontext)
	cmdcontext = fmt.Sprintf("ansible-playbook -i infra-om-%s/hosts infra-om-%s/configure-cmd-binaries.yml", v1, v1)
	command_run_check(cmdcontext)
	log.Println(cmdcontext)
	cmdcontext = fmt.Sprintf("ansible-playbook -i infra-om-%s/hosts  infra-om-%s/configure-tenant-clusterroles.yml -e '{test_env: 5G}'", v1, v1)
	command_run_check(cmdcontext)
	log.Println(cmdcontext)
}

func python_env() {
	v1 := utils.GlobalObject.ScriptsVersion
	log.Println("--------start perpare OE python env--------")
	set_python_info := "\"8i \\ \\ \\ \\ ssh \\$SSH_PARAMS \\$node \\\"sudo ln -s -f /usr/libexec/platform-python3.6 /usr/local/bin/python\\\"\""
	cmdcontext := "sed -i " + set_python_info + fmt.Sprintf(" infra-om-%s/roles/configure_python/templates/python.sh.jinja", v1)
	log.Println(cmdcontext)
	command_run_check(cmdcontext)

	cmdcontext = fmt.Sprintf("ansible-playbook -i infra-om-%s/hosts infra-om-%s/configure-python.yml", v1, v1)
	command_run_check(cmdcontext)
}

func python_env_other() {
	log.Println("--------start login OE and install requests, pyyaml--------")
	set_proxy_info1 := "export https_proxy=http://10.158.100.3:8080"
	set_proxy_info2 := "export http_proxy=http://10.158.100.3:8080"
	set_auth_info1 := "sudo chmod -R 755 /usr/local/lib "
	set_auth_info2 := "sudo chmod -R 755 /usr/local/bin/"
	set_auth_info3 := "sudo mkdir -p /usr/local/lib64 && sudo chmod -R 755 /usr/local/lib64"

	sshlib.ClientServe(set_auth_info1)
	out, _ := sshlib.ClientServe("ls -al /usr/local/lib")
	log.Println(string(out))
	sshlib.ClientServe(set_auth_info2)
	out, _ = sshlib.ClientServe("ls -al /usr/local/bin/")
	log.Println(string(out))
	sshlib.ClientServe(set_auth_info3)
	out, _ = sshlib.ClientServe("ls -al /usr/local/lib64")
	log.Println(string(out))

	pip_info := "python -m pip install requests && python -m pip install pyyaml==5.4.1"
	cmdcontext := set_proxy_info1 + " && " + set_proxy_info2 + " && " + pip_info
	log.Println(cmdcontext)

	buf, err := sshlib.ClientServe(cmdcontext)
	if err != nil {
		log.Println("sshlib.ClientServe err: ", err)
		return
	}
	log.Println(string(buf))
}
