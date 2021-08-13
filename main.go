// Copyright (c) 2015-2021 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"os"

	"github.com/minio/cli"
)

var helpTemplate = `NAME:
{{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} [arguments...]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLE:
	{{.HelpName}} https://endpoint ACCESSKEY SECRETKEY BUCKETNAME
	
`

var createCmd = cli.Command{
	Name:               "create",
	Usage:              "creates a dataset on the specified bucket",
	Action:             mainCreate,
	Flags:              insecureFlag,
	CustomHelpTemplate: helpTemplate,
}

var verifyCmd = cli.Command{
	Name:   "verify",
	Usage:  "Downloads and verifies the hash of the objects created by the create command",
	Action: mainVerify,
	Flags:  insecureFlag,
	//CustomHelpTemplate: helpTemplate,
}

var deleteCmd = cli.Command{
	Name:   "delete",
	Usage:  "deletes all the objects in the specified bucket",
	Action: mainDelete,
	//CustomHelpTemplate: helpTemplate,
}

//list of commands
var appCmds = []cli.Command{
	createCmd,
	verifyCmd,
	deleteCmd,
}

//flags that are used
var (
	insecureFlag = []cli.Flag{
		cli.BoolFlag{
			Name:  "insecure",
			Usage: "disable SSL certificate verification",
		},
	}
)

func main() {
	app := cli.NewApp()
	app.UsageText = "A tool to validate the consistency of data uploaded to MinIO."
	app.Commands = appCmds
	app.Action = func(ctx *cli.Context) error {
		if ctx.Args().First() == "" {
			cli.ShowAppHelp(ctx)
		}

		return nil
	}
	app.Usage = "A tool to validate the consistency of data uploaded to MinIO."
	app.Author = "MinIO, Inc."
	app.HideHelpCommand = true
	app.Flags = append(insecureFlag)
	app.Run(os.Args)
}
