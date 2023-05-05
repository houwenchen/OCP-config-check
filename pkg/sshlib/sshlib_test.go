package sshlib

import (
	"log"
	"testing"
)

func TestClientServe(t *testing.T) {
	buf, err := ClientServe("ls -al")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(buf))
}

func TestClientDownloadFile(t *testing.T) {
	remoteFile := "/var/home/core/tmp/openshift-ptp.json"
	localFile := "C:/houwenchen/work/go/OCP_with_Command/pkg/sshlib/testData/openshift-ptp.json"
	_, err := ClientDownloadFile(remoteFile, localFile)
	if err != nil {
		log.Println(err)
	}
}

func TestClientClientUploadFile(t *testing.T) {
	remoteFile := "/var/home/core/tmp/testfile.txt"
	localFile := "C:/houwenchen/work/go/OCP_with_Command/pkg/sshlib/testData/testfile.txt"
	_, err := ClientUploadFile(localFile, remoteFile)
	if err != nil {
		log.Println(err)
	}
}
