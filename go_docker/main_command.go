package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"go_docker/container"
	"go_docker/cgroups/subsystems"
	"os"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit go_docker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name: "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name: "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name: "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name: "v",
			Usage: "volume",
		},
		cli.BoolFlag{
			Name: "d",
			Usage: "detach container",
		},
		cli.BoolFlag{
			Name: "l",
			Usage: "log file container",
		},
		cli.StringFlag{
			Name: "name",
			Usage: "container name",
		},
		cli.StringSliceFlag{
			Name: "e",
			Usage: "set environment",
		},
	},

	Action: func(context *cli.Context) error {
		if len(context.Args()) <1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _,arg := range context.Args(){
			cmdArray = append(cmdArray, arg)
		}
		//cmd := context.Args().Get(0)
		tty := context.Bool("it")
		detach := context.Bool("d")

		//get image name
		imageName := cmdArray[0]
		cmdArray = cmdArray[1:]

		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet: context.String("cpuset"),
			CpuShare: context.String("cpushare"),
		}

		if detach && tty {
			return fmt.Errorf("it and d parameter can not both provided")
		}

		volume := context.String("v")
		logfile := context.Bool("l")

		log.Infof("createTty %v", tty)
		containerName := context.String("name")

		envSlice := context.StringSlice("e")
		Run(tty, cmdArray, resConf, volume,containerName,logfile,imageName,envSlice)
		return nil
	},
}

var initCommand = cli.Command{
	Name: "init",
	Usage: "Init container process run user's process in container.Do not call it outside",

	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		//cmd := context.Args().Get(0)
		//log.Infof("command %s", cmd)
		//err := container.RunContainerInitProcess(cmd,nil)
		err := container.RunContainerInitProcess()
		return err
	},
}

var commitCommand = cli.Command{
	Name: "commit",
	Usage: "commit a container into image",
	Action: func(context *cli.Context) error{
		if len(context.Args()) < 2 {
			return fmt.Errorf("Missing image name and container name!")
		}
		containerName := context.Args().Get(0)
		imageName := context.Args().Get(1)
		commitContainer(containerName,imageName)
		return nil
	},
}


var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list all the containers",
	Action: func(context *cli.Context) error {
		ListContainers()
		return nil
	},
}


var logCommand = cli.Command{
	Name: "logs",
	Usage: "print logs of a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Please input your container name")
		}
		containerName := context.Args().Get(0)
		logContainer(containerName)
		return nil
	},
}


var execCommand = cli.Command{
	Name: "exec",
	Usage: "exec a command into container",
	Action: func(context *cli.Context) error {
		//This is for callback
		if os.Getenv(ENV_EXEC_PID) != "" {
			log.Infof("pid callback pid %s", os.Getgid())
			return nil
		}

		// exec containername sh
		if len(context.Args()) < 2 {
			return fmt.Errorf("Missing container name or command")
		}
		containerName := context.Args().Get(0)
		var commandArray []string
		for _, arg := range context.Args().Tail() {
			commandArray = append(commandArray, arg)
		}
		ExecContainer(containerName, commandArray)
		return nil
	},
}


var stopCommand = cli.Command{
	Name: "stop",
	Usage: "stop a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
		StopContainer(containerName)
		return nil
	},
}


var removeCommand = cli.Command{
	Name: "rm",
	Usage: "remove unused containers",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
		RemoveContainer(containerName)
		return nil
	},
}