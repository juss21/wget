package wget

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	//"os/exec"
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
	fmt.Print("Connecting ", url_split[2], " ("+url_split[2]+")|", add[0], "|:"+Port+"...")
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
	fmt.Println()
	size, filetype := FileInfo(url_split[4], urls)
	fmt.Println("Length:", size, "[" + filetype + "]")
	fmt.Println("Saving to:", url_split[4])
	fmt.Println()
	fmt.Println(url_split[4]+ "\t\t\t\t" + "s")
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

func FileInfo(FileName, url string) (size int64, FileType string) {
	tr := new(http.Transport)
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
		return 0, ""
	}

	size = resp.ContentLength
	bytes, _ := ioutil.ReadAll(resp.Body)
	contentType := http.DetectContentType(bytes)

	return size, contentType
}
