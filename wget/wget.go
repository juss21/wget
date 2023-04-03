package wget

import (
	"bufio"
	//"bytes"
	"io"
	//"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/juju/ratelimit"
	"github.com/schollz/progressbar/v3"
)

// Make the GET request to a url, return the response
func getResponse(link, httpmethod, shorturl, fileName, givenpath string) (*http.Response, string, string) {
	net.LookupPort("tcp", httpmethod)
	ip, errx := net.LookupIP(shorturl)
	errorHandler(errx, true)

	u, erro := url.Parse(link)
	errorHandler(erro, true)
	Port := GetPortFromHttpMethod(u.Scheme)

	// resolving
	if len(ip) > 1 {
		doLogging("Resolving "+shorturl+" ("+shorturl+")..."+ip[0].String()+","+ip[1].String(), true)
	} else {
		doLogging("Resolving "+shorturl+" ("+shorturl+")..."+ip[0].String(), true)
	}

	doLogging("Connecting "+shorturl+" ("+shorturl+")|"+ip[0].String()+"|:"+Port+"...", false)

	size, filetype := FileInfo(fileName, link)
	if !strings.Contains(fileName, ".") {
		Type := strings.Split(filetype, "/")
		fileName += "." + strings.TrimSuffix(Type[1], "]")
	}
	//client
	req, err := http.NewRequest("GET", link, nil)
	errorHandler(err, true)

	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, err3 := http.DefaultClient.Do(req)
	errorHandler(err3, true)
	doLogging(" connected.", true)

	doLogging(httpmethod+"request sent, awaiting response... "+resp.Status, true)
	if resp.StatusCode != 200 {
		doLogging("Location: "+link+" [following]", true)
		getResponse(link, httpmethod, shorturl, fileName, givenpath)
	}

	doLogging("", true)

	a := strconv.FormatInt(size, 10)
	doLogging("Length:"+a+CalcSize(size)+"["+filetype+"]", true)

	doLogging("Saving to file: "+fileName, true)
	doLogging("", true)

	return resp, a, fileName
}

// Write the response of the GET request to file
func writeToFile(directory, fileName string, resp *http.Response) (elapsed time.Duration, data int64) {
	var file *os.File
	if Flags.P_Flag == "" {
		createPath("downloads/")
		if strings.HasPrefix(directory, "downloads/") {
			createPath(directory)
		} else {
			createPath("downloads/" + directory)
		}
		fileo, erro := os.OpenFile(directory+"/"+fileName, os.O_CREATE|os.O_WRONLY, 0o644)
		errorHandler(erro, false)
		file = fileo
	} else {
		createPath(directory)
		fileo, erro := os.OpenFile(directory+fileName, os.O_CREATE|os.O_WRONLY, 0o644)
		errorHandler(erro, false)
		file = fileo
	}

	start := time.Now()
	defer file.Close()

	bufferedWriter := bufio.NewWriterSize(file, bufSize)
	if Flags.B_Flag {
		if Flags.RateLimit_Flag != "" {
			limit := ConvertLimit(strings.ToLower(Flags.RateLimit_Flag)) 
			bucket := ratelimit.NewBucketWithRate(float64(limit), int64(limit))
			data, _ = io.Copy(bufferedWriter, ratelimit.Reader(resp.Body, bucket))
		} else {
			data, _ = io.Copy(bufferedWriter, resp.Body)
		}
	} else {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"Downloading: "+fileName,
		)
		if Flags.RateLimit_Flag != "" {
			limit := ConvertLimit(strings.ToLower(Flags.RateLimit_Flag))
			bucket := ratelimit.NewBucketWithRate(float64(limit), int64(limit))
			body := ratelimit.Reader(resp.Body, bucket)
			time.Sleep(1 * time.Second)
			data, _ = io.Copy(io.MultiWriter(bufferedWriter, bar), body)
		} else {
			data, _ = io.Copy(io.MultiWriter(bufferedWriter, bar), resp.Body)
		}
	}
	t := time.Now()
	elapsed = t.Sub(start)

	return elapsed, data
}

func FileInfo(FileName, url string) (size int64, FileType string) {
	req, err := http.NewRequest("GET", url, nil)
	errorHandler(err, true)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, err3 := http.DefaultClient.Do(req)
	errorHandler(err3, true)

	size = resp.ContentLength

	buf := make([]byte, 512)

	var contentType string
	_, err2 := resp.Body.Read(buf)
	if err2 != nil {
		contentType = "Couldn't read file"
		return size, contentType
	}
	contentType = http.DetectContentType(buf)
	return size, contentType
}
