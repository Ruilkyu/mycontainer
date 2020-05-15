package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go_docker/container"
	"go_docker/cgroups/subsystems"
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
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet: context.String("cpuset"),
			CpuShare: context.String("cpushare"),
		}

		volume := context.String("v")
		Run(tty, cmdArray, resConf, volume)
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
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing image name!")
		}
		imageName := context.Args().Get(0)
		commitContainer(imageName)
		return nil
	},
}