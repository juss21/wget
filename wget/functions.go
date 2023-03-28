package wget

import (
	"fmt"
	"net"
	"strconv"

	//"io/ioutil"

	"net/http"
	"net/url"
	"os"
	"strings"
	//"os/exec"
)

// Check if we received an error on our last function call
func errorChecker(err error) {
	if err != nil {
		doLogging(err.Error(), true)
		os.Exit(0)
	}
}

// Read data from url
func sliceUrl(url string) (rurl, cleanurl, givenfilename, givenpath, httpmethod string) {
	var rebuilt []string
	fmt.Println(url)
	split := strings.Split(url, "/")

	for i := 0; i < len(split); i++ {
		if split[i] != "" {
			rebuilt = append(rebuilt, split[i])
		}
	}

	rurl = url
	cleanurl = rebuilt[1]
	givenfilename = rebuilt[2]
	if len(split) != 4 {
		givenpath = rebuilt[2]
		givenfilename = rebuilt[3]
	}
	httpmethod = rebuilt[0]

	return rurl, cleanurl, givenfilename, givenpath, httpmethod
}

// Make the GET request to a url, return the response
func getResponse(link, httpmethod, shorturl string) *http.Response {
	net.LookupPort("tcp", httpmethod)

	u, err := url.Parse(link)
	Port := GetPort(u.Scheme)

	doLogging("Resolving "+shorturl+" ("+shorturl+")... "+httpmethod+" "+shorturl, true)
	tr := new(http.Transport)
	doLogging("Connecting "+shorturl+" ("+shorturl+")|"+httpmethod+"|:"+Port+"...", false)

	client := &http.Client{Transport: tr}
	resp, err := client.Get(link)
	errorChecker(err)
	doLogging(" connected.", true)
	doLogging("HTTP request sent, awaiting response... "+resp.Status, true)
	if resp.StatusCode != 200 {
		doLogging("Location: "+link+" [following]", true)
		getResponse(link, httpmethod, shorturl)
	}

	doLogging("", true)
	tempFile := ""
	if Flags.O_Flag != "" {
		tempFile = Flags.O_Flag
	}

	size, filetype := FileInfo(tempFile, link)
	a := strconv.FormatInt(size, 10)
	doLogging("Length:"+a+CalcSize(size)+"["+filetype+"]", true)
	doLogging("Saving to:"+tempFile, true)
	doLogging("", true)
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
func CalcSize(numb int64) string {
	var output string

	numstr := strconv.FormatInt(numb, 10)
	num, _ := strconv.Atoi(numstr)

	if num < 1024 {
		output = "(" + strconv.Itoa(num) + "B)"
	} else if num >= 1024 && num < 1048576 {
		output = "(" + strconv.Itoa(num/1024) + "K)"
	} else if num >= 1048576 && num < 134217728 {
		output = "(" + strconv.Itoa(num/1048576) + "M)"
	} else {
		output = "(" + strconv.Itoa(num/1073700000) + "G)"
	}

	return output
}

func FileInfo(FileName, url string) (size int64, FileType string) {
	tr := new(http.Transport)
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	errorChecker(err)
	size = resp.ContentLength
	buf := make([]byte, 512)
	_, err2 := resp.Body.Read(buf)
	errorChecker(err2)
	contentType := http.DetectContentType(buf)
	return size, contentType
}
