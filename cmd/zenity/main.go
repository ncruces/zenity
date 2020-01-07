package main

import (
	"os"

	"github.com/ncruces/zenity"
	"github.com/ncruces/zenity/internal/cmd"
)

func main() {
	cmd.Command = true

	file, err := zenity.SelectFile()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(255)
	}
	if file == "" {
		os.Exit(1)
	}
	os.Stdout.WriteString(file)
	os.Stdout.WriteString("\n")
}
