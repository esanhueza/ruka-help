package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
)

var flags []cli.Flag

func init() {

	flags = []cli.Flag{}
}

func main() {

	app := cli.NewApp()
	app.Usage = "Chat"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Erick Sanhueza",
			Email: "esanhueza@zohomail.com",
		},
	}
	app.Flags = flags

	app.Action = StartListener

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
