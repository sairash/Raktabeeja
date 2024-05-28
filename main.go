package main

import (
	"log"
	"os"
	"raktabeeja/files"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "RaktaBeeja",
		Usage: "Decentralize deez files",
		Commands: []*cli.Command{
			{
				Name:    "export",
				Aliases: []string{"e"},
				Usage:   "Exports the file to Chunks",
				Action: func(cCtx *cli.Context) error {
					files.ExportFile(cCtx.Args().First())
					return nil
				},
			},
			{
				Name:    "read",
				Aliases: []string{"r"},
				Usage:   "Reads the file from the Chunk",
				Action: func(cCtx *cli.Context) error {
					files.ReadByt(cCtx.Args().First())
					return nil
				},
			},
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "Test the file from the Chunk",
				Action: func(cCtx *cli.Context) error {
					files.ReadByt(cCtx.Args().First())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// name node
// master node
// data node
// client node
