package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const usage = `learn docker and go!`

func main(){
	app := cli.NewApp()
	app.Name = "go_docker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
		commitCommand,
		listCommand,
		logCommand,
	}

	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
