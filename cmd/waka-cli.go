package main

import (
	"fmt"
	"os"

	"github.com/infra-whizz/waka"

	"github.com/isbm/go-nanoconf"
	"github.com/urfave/cli/v2"
)

// Application dispatcher
func appDispatcher(ctx *cli.Context) error {
	if ctx.String("schema") == "" {
		return fmt.Errorf("Error: Image schema path was not provided. Try --help, perhaps?\n")
	}

	waka.NewWaka().LoadSchema(ctx.String("schema")).Build()

	cli.ShowAppHelpAndExit(ctx, 0)
	return nil
}

func main() {
	appname := "waka"
	confpath := nanoconf.NewNanoconfFinder(appname).DefaultSetup(nil)

	app := &cli.App{
		Version: "0.1 Alpha",
		Name:    appname,
		Usage:   "image builder with config mgmt powers",
		Action:  appDispatcher,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Path to the configuration file",
				Required: false,
				Value:    confpath.SetDefaultConfig(confpath.FindFirst()).FindDefault(),
			},
			&cli.StringFlag{
				Name:     "schema",
				Aliases:  []string{"s"},
				Usage:    "Path to the image schema",
				Required: false,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
