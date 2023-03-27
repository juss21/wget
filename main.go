package main

import (
	"fmt"
	"os"
	"wget/wget"
)

//flags
/*
-B



*/

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Incorrect usage!")
		printUsage()
	}

	wget.Run(os.Args[1], "") // wget
}

func printUsage() {
	fmt.Println("wget: missing URL")
	fmt.Println("Usage: go run . [OPTION]... [URL]...")
	fmt.Println()
	fmt.Println("Try `go run  . --help' for more options.")
	os.Exit(0)
}

func printHelp() {
	fmt.Println("USAGE:")

	fmt.Println("go run . <download_link>")
	fmt.Println("EXAMPLE:")
	fmt.Println("go run . https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg")
	os.Exit(0)
}
