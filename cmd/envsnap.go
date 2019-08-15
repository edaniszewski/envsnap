package main

import (
	"fmt"
	"os"

	"github.com/edaniszewski/envsnap/pkg"
)

func main() {
	app := pkg.NewApp()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
