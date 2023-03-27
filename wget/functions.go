package wget

import (
	"fmt"
	//"io"
	"net"
	"net/url"
	"net/http"
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
	host, port, err := net.SplitHostPort(u.Host)
	fmt.Println(u.Host, urls)
	fmt.Println(host, port)

	fmt.Println("Resolving", url_split[2], "("+url_split[2]+")...", add[0], add[1])
	tr := new(http.Transport)
	fmt.Println("Connecting", url_split[2], "("+url_split[2]+") |"+string(add[0])+"|")
	client := &http.Client{Transport: tr}
	resp, err := client.Get(urls)
	//fmt.Println(url_split[2])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("HTTP request sent, awaiting response... ", resp.StatusCode)
	if resp.StatusCode != 200 {
		return resp
	}
	errorChecker(err)
	fmt.Println("")
	return resp
}
