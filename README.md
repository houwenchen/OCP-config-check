功能：
1. 新安装集群检查：
    检查mtu、driver version、network interface num、ptp同步即失败后恢复；
2. vDU、vCU部署前环境配置：
    包括 https://confluence.ext.net.nokia.com/display/ASAiCBTS/Step0_5+Tenant+configuration+before+VCU+and+VDU+deploy+after+OCP+post+configuration 的所有步骤。

使用方法：
1. easytool客户端。


修改记录：
2022-12-11：
1. 增加modify_hosts函数，实现修改infra-om-2.26.0/hosts文件；
2. 增加mlog模块，实现log打印到控制台的同时，保存log到log.txt文件。
2022-12-15: 
1. 补充sshlib模块，实现sftp上传、下载文件；
2. 增加check_ptp函数，实现OE20同步检查，同步接口修改。
2023-2-8:
1. 增加命令行功能，起名为easytool。
2023-2-13:
1. 使用C-S架构，
2. eg：
     #参数相关命令
      easytool para get
      easytool para set  #写入全局变量中，服务器关闭会加载默认配置，有很多flag，可以用-h命令获取。

    #check命令
      easytool check all
      easytool check capacity
      easytool check clusternode
      easytool check driver
      easytool check mtu
      easytool check pyenv
      easytool check ptp

    #config命令
      easytool config all
      easytool config tenantuser
      easytool config py
2023-2-14:
1. 添加build脚本。


备注：
运行环境：
1. go1.17.8;
2. 将本机的GOPATH下的pkg下的库文件拷贝到运行环境的GOPATH下的pkg。---弃用

TODO：后面尝试打包这个项目，方便在没有go的环境中运行。---可以使用vendor这种依赖解决方法--DONE
TODO: 运行环境如果有k8s客户端的话，就可以不用ssh库，可以用k8s的api来做一些操作。
TODO：增加build脚本--DONE
TODO：config模块优化，增加细节输出
TODO：客户端log模块优化
TODO：后续考虑下怎么做成多线程
TODO：看看ansible-playbook怎么安装的，有没有办法不用python依赖
TODO: 增加不同版本ocp使用的脚本，不同的环境使用不同的文件，这样在多个账户的时候可以不冲突，文件夹可以设置为：
         