package sshlib

import (
	"fmt"
	"io"
	"log"
	_ "ocp-check-config/pkg/mlog"
	"ocp-check-config/pkg/utils"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	Ip         string
	Port       int
	User       string
	Passwd     string
	SClient    *ssh.Client
	SFClient   *sftp.Client
	SshSession *ssh.Session
}

func NewSshClient() *SshClient {
	return &SshClient{
		Ip:     utils.GlobalObject.SshIP,
		Port:   22,
		User:   utils.GlobalObject.User,
		Passwd: utils.GlobalObject.Passwd,
	}
}

// according global config create ssh client
func (scli *SshClient) SetSClient() error {
	//load global config
	config := &ssh.ClientConfig{
		User:            scli.User,
		Auth:            []ssh.AuthMethod{ssh.Password(scli.Passwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	addr := fmt.Sprintf("%v:22", scli.Ip)

	//create ssh client
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("unable to create ssh conn, err: ", err)
		return err
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Fatal("unable to create sftp conn, err: ", err)
		return err
	}
	scli.SClient = sshClient
	scli.SFClient = sftpClient
	return nil
}

// create ssh session
func (scli *SshClient) SetSession() error {
	sshSession, err := scli.SClient.NewSession()
	if err != nil {
		log.Fatal("unable to create ssh session, err: ", err)
		return err
	}
	scli.SshSession = sshSession
	return nil
}

// send command and return stdout
func (scli *SshClient) RunCommand(cmd string) ([]byte, error) {
	buf, err := scli.SshSession.CombinedOutput(cmd)
	if err != nil {
		log.Fatal("scli.SshSession.CombinedOutput err: ", err)
		return []byte{}, err
	}
	return buf, nil
}

// downfile
func (scli *SshClient) DownloadFile(remoteFile, localFile string) (int, error) {
	source, err := scli.SFClient.Open(remoteFile)
	if err != nil {
		log.Fatal("unable to open remote file")
		return -1, err
	}
	defer source.Close()

	target, err := os.OpenFile(localFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("unable to open local file")
		return -1, err
	}
	defer target.Close()

	n, err := io.Copy(target, source)
	if err != nil {
		log.Fatal("copy file error")
		return -1, err
	}
	return int(n), nil
}

// upload file, local open file ---> read ---> remote create file ---> write
func (scli *SshClient) UploadFile(localFile, remoteFile string) (int, error) {
	file, err := os.Open(localFile)
	if nil != err {
		log.Fatal("unable to open local file")
		return -1, err
	}
	defer file.Close()

	ftpFile, err := scli.SFClient.Create(remoteFile)
	if nil != err {
		log.Fatal("unable to create file in remote host")
		return -1, err
	}
	defer ftpFile.Close()

	fileByte, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("ioutil.ReadAll err")
		return -1, err
	}

	n, err := ftpFile.Write(fileByte)
	if err != nil {
		log.Fatal("ftpFile.Write err")
		return -1, err
	}
	return n, nil
}

func ClientServe(cmd string) ([]byte, error) {
	cli := NewSshClient()
	cli.SetSClient()
	cli.SetSession()
	// defer cli.SshSession.Close()

	buf, err := cli.RunCommand(cmd)
	return buf, err
}

func ClientDownloadFile(remoteFile, localFile string) (int, error) {
	cli := NewSshClient()
	cli.SetSClient()
	cli.SetSession()
	// defer cli.SshSession.Close()

	n, err := cli.DownloadFile(remoteFile, localFile)
	return n, err
}

func ClientUploadFile(localFile, remoteFile string) (int, error) {
	cli := NewSshClient()
	cli.SetSClient()
	cli.SetSession()
	// defer cli.SshSession.Close()

	n, err := cli.UploadFile(localFile, remoteFile)
	return n, err
}
