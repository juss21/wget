package wget

import (
	//"fmt"
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
	//fmt.Println(url)
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
func getResponse(link, httpmethod, shorturl, fileName string) (*http.Response, string) {
	net.LookupPort("tcp", httpmethod)
	ip, errx := net.LookupIP(shorturl)
	errorChecker(errx)

	u, err := url.Parse(link)
	errorChecker(err)
	Port := GetPort(u.Scheme)

	doLogging("Resolving "+shorturl+" ("+shorturl+")... "+(ip[0]).String()+ ", "+(ip[1]).String(), true)
	//tr := new(http.Transport)
	doLogging("Connecting "+shorturl+" ("+shorturl+")|"+(ip[0]).String()+"|:"+Port+"...", false)

	//client := &http.Client{Transport: tr}
	req, err2 := http.NewRequest("GET", link, nil)
	errorChecker(err2)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, err3 := http.DefaultClient.Do(req)
	//resp, err := client.Get(link)
	errorChecker(err3)
	doLogging(" connected.", true)
	doLogging(httpmethod+"request sent, awaiting response... "+resp.Status, true)
	if resp.StatusCode != 200 {
		doLogging("Location: "+link+" [following]", true)
		getResponse(link, httpmethod, shorturl, fileName)
	}

	doLogging("", true)

	size, filetype := FileInfo(fileName, link)
	a := strconv.FormatInt(size, 10)
	doLogging("Length:"+a+CalcSize(size)+"["+filetype+"]", true)

	doLogging("Saving to: "+fileName, true)
	doLogging("", true)
	return resp, a
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
	num, err := strconv.Atoi(numstr)
	errorChecker(err)
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
	req, err := http.NewRequest("GET", url, nil)
	errorChecker(err)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, err3 := http.DefaultClient.Do(req)
	errorChecker(err3)

	size = resp.ContentLength

	buf := make([]byte, 512)

	var contentType string
	_, err2 := resp.Body.Read(buf)
	if err2 != nil {
		contentType = "Couldn't read file"
		goto exit
	}

	contentType = http.DetectContentType(buf)
exit:
	return size, contentType
}
