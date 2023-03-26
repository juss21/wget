package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Incorrect usage!")
		printUsage()
	}
	app.wget(os.Args[2]) // wget
}

func printUsage() {
	fmt.Println("USAGE:")
	fmt.Println("go run . <download_link>")
	fmt.Println("EXAMPLE:")
	fmt.Println("go run . https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg")
	os.Exit(0)
}
