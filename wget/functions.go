package wget

import (
	"fmt"
	"strings"
	//"io"
	"net"
	"net/http"
	"net/url"
	"os"
)

// Check if we received an error on our last function call
func errorChecker(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// Make the GET request to a url, return the response
func getResponse(urls string, url_split []string) *http.Response {
	add, err := net.LookupIP(url_split[2])

	net.LookupPort("tcp", "https")

	u, err := url.Parse(urls)
	Port := GetPort(u.Scheme)

	fmt.Println("Resolving", url_split[2], "("+url_split[2]+")...", add[0], add[1])
	tr := new(http.Transport)
	fmt.Print("Connecting ", url_split[2], " ("+url_split[2]+")|", add[0], "|:" + Port + "...")
	client := &http.Client{Transport: tr}
	resp, err := client.Get(urls)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(" connected.")
	fmt.Print("HTTP request sent, awaiting response... ", resp.Status)
	if resp.StatusCode != 200 {
		return resp
	}
	errorChecker(err)
	fmt.Println("")
	return resp
}

func GetPort(s string) (port string) {
	if strings.ToLower(s) == "https" {
		port = "443"
	} else if strings.ToLower(s) == "http" {
		port = "80"
	} else if strings.ToLower(s) == "telnet" {
		port = "23"
	} else if strings.ToLower(s) == "ftp" {
		port = "21"
	}
	return port
}

func FileInfo() (size int64, FileType string){
	fi, err := os.Stat("/path/to/file")
	if err != nil {
		return 0, ""
	}
	// get the size
	size = fi.Size()
	return size, "img"
}