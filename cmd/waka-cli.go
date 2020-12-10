package main

import (
	"fmt"
	"os"

	"github.com/infra-whizz/waka"
	"github.com/isbm/go-nanoconf"
	"github.com/urfave/cli/v2"
)

// Description builder
func descriptionBuilder(ctx *cli.Context) error {
	cli.ShowAppHelpAndExit(ctx, 1)
	return nil
}

// Image builder
func imageBuilder(ctx *cli.Context) error {
	if ctx.String("schema") == "" {
		return fmt.Errorf("Error: Image schema path was not provided. Try --help, perhaps?\n")
	}
	waka.NewWaka().
		SetCleanupOnExit(!ctx.Bool("debug")).
		SetBuildOutput(ctx.String("output")).
		LoadSchema(ctx.String("schema")).
		Build(ctx.Bool("force"))

	return nil
}

func main() {
	appname := "waka"
	confpath := nanoconf.NewNanoconfFinder(appname).DefaultSetup(nil)

	app := &cli.App{
		Version: "0.1 Alpha",
		Name:    appname,
		Usage:   "image builder with config mgmt powers",
	}

	app.Commands = []*cli.Command{
		{
			Name:   "decription",
			Usage:  "Create a description with collection of required modules",
			Action: descriptionBuilder,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "config",
					Aliases:  []string{"c"},
					Usage:    "Path to the description configuration file",
					Required: false,
				},
			},
		},
		{
			Name:   "image",
			Usage:  "Build an image, based on an existing description",
			Action: imageBuilder,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "config",
					Aliases:  []string{"c"},
					Usage:    "Path to the image layout configuration file",
					Required: false,
					Value:    confpath.SetDefaultConfig(confpath.FindFirst()).FindDefault(),
				},
				&cli.StringFlag{
					Name:     "schema",
					Aliases:  []string{"s"},
					Usage:    "Path to the image schema",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "output",
					Aliases:  []string{"o"},
					Usage:    "Path to build output (default is $SCHEMA/build)",
					Required: false,
				},
				&cli.BoolFlag{
					Name:     "debug",
					Aliases:  []string{"d"},
					Usage:    "Leave mounts untouched, log debug messages",
					Value:    false,
					Required: false,
				},
				&cli.BoolFlag{
					Name:     "force",
					Aliases:  []string{"f"},
					Usage:    "Flush previous builds",
					Required: false,
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
