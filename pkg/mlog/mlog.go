package mlog

import (
	"fmt"
	"io"
	"log"
	"os"
)

/*
	定制化log模块，将log同时输出文件和控制台
*/

func init() {
	file_path := "log.txt" //在本文件夹下保存log

	//检查是否存在以前的log，存在的时候先删除文件
	if _, err := os.Stat(file_path); err == nil {
		err = os.Remove(file_path)
		if err != nil {
			panic(fmt.Sprintf("os.Remove file error: %v", err))
		}
	}

	//创建文件，这里其实也可以直接创建
	logFile, err := os.OpenFile(file_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(fmt.Sprintf("os.OpenFile error: %v", err))
	}

	//保证log同时输出到控制台和文件
	f := io.MultiWriter(logFile, os.Stdout)

	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
