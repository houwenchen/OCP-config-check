package main

import (
	"log"
	"ocp-check-config/pkg"
	_ "ocp-check-config/pkg/mlog"
)

func main() {
	// cmd.Execute()
	app := pkg.NewApp()
	err := app.Run()
	if err != nil {
		log.Println("app run err: ", err)
	}
}
