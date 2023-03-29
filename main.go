package main

import (
	"fmt"
	"os"
	"wget/wget"
)

func main() {
	wget.BuildFlags() // build flags
	if len(os.Args) == 1 {
		printUsage()
	} else if wget.Flags.H_Flag {
		printHelp()
	}

	wget.Run() // run wget
}

func printUsage() {
	fmt.Println("wget: missing URL")
	fmt.Println("Usage: go run . [OPTION]... [URL]...")
	fmt.Println()
	fmt.Println("Try `go run  . --help' for more options.")
	os.Exit(0)
}

func printHelp() {
	fmt.Print("USAGE: ")
	fmt.Println("go run . <download_link>")

	fmt.Print("EXAMPLE: ")
	fmt.Println("go run . https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg")
	os.Exit(0)
}
