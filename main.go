package main

import (
	"os"

	"github.com/jakubknejzlik/dns-deploy/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "dns-deploy"
	app.Usage = "Deploy DNS configuration to various DNS providers."
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		cmd.RunCommand(),
	}

	app.Run(os.Args)
}
