package main

import (
	"fmt"
	"os"
	"wget/wget"
)

/*todo ja gtn
todo: luua wget-log (-B flag)
todo: teha mirror	(-mirror flag)
todo: rate-limiter	(--rate-limit flag)

gtn: average dl etc (numbrid overall) katki nüüd doLogging() funci pärast :)
gtn: downloads.txt
on 2 tekstifaili mis viskavad errorisse
HTTP request sent, awaiting response... 403 Forbidden
see on sellepärast, et see veebileht ei lase tõmmata sealt
peame oma wgetile lisama browseri headeri vms request.Header.Add()

todo: code cleanup
muidu redis vist, paneb auditi käima v?
*/

func main() {
	wget.BuildFlags()
	if len(os.Args) == 1 {
		printUsage()
	} else if wget.Flags.H_Flag {
		printHelp()
	}

	wget.Run() // wget
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
