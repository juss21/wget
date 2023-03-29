package wget

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/juju/ratelimit"
	"github.com/schollz/progressbar/v3"
)

const bufSize = 1024 * 8

var wg sync.WaitGroup

func Run() {
	if Flags.B_Flag {
		os.Truncate("wget-log", 0)
		fmt.Println("Output will be written to \"wget-log\".")
	}

	// looping all the links saved in Flags
	for i, file := range Flags.Links {
		wg.Add(1)

		url, shorturl, givenfilename, givenpath, httpmethod := sliceUrl(file)
		tempFile := givenfilename
		if Flags.O_Flag != "" {
			tempFile = Flags.O_Flag
		}

		download_started := time.Now().Format("--2006-01-02 15:04:05--")
		doLogging(download_started+"\t"+url, true)
		go startDownload(url, shorturl, tempFile, givenpath, httpmethod)
		if i >= 2 {
			break
		}
	}
	wg.Wait()

}

func startDownload(url, shorturl, filename, givenpath, httpMethod string) {

	if Flags.P_Flag != "" {
		givenpath = Folder(Flags.P_Flag)
	}

	response, filesize := getResponse(url, httpMethod, shorturl, filename)
	elapsed, data := writeToFile(givenpath, filename, response)

	h, _ := time.ParseDuration(elapsed.String())
	size, _ := strconv.Atoi(filesize)
	var AvgDown float64
	if h.Seconds() < 1 {
		AvgDown = float64(size) * (h.Seconds()) / 10
	} else {
		AvgDown = float64(size) / (h.Seconds()) / 1000000
	}
	DownloadedDataInt := strconv.FormatInt(data, 10)

	doLogging(time.Now().Format("2006-01-02 15:04:05 - Download completed! "), false)
	doLogging("["+string(DownloadedDataInt)+"/"+filesize+"]", true)
	doLogging("Time elapsed: ", false)
	doLogging(shortTimeDur(elapsed), false)
	doLogging(" Average download speed: ", false)
	math := math.Round(AvgDown*10) / 10
	conv := strconv.FormatFloat(math, 'f', -1, 64)
	doLogging(conv, false)
	doLogging(" MB/s"+"\nFile location: '"+givenpath+"/"+filename+"'", true)

	wg.Done()
}

// Make the GET request to a url, return the response
func getResponse(link, httpmethod, shorturl, fileName string) (*http.Response, string) {
	net.LookupPort("tcp", httpmethod)
	ip, errx := net.LookupIP(shorturl)
	errorHandler(errx, true)

	u, err := url.Parse(link)
	errorHandler(err, true)
	Port := GetPortFromHttpMethod(u.Scheme)

	// resolving
	if len(ip) > 1 {
		doLogging("Resolving "+shorturl+" ("+shorturl+")..."+ip[0].String()+","+ip[1].String(), true)
	} else {
		doLogging("Resolving "+shorturl+" ("+shorturl+")..."+ip[0].String(), true)
	}

	doLogging("Connecting "+shorturl+" ("+shorturl+")|"+ip[0].String()+"|:"+Port+"...", false)

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
		getResponse(link, httpmethod, shorturl, fileName)
	}

	doLogging("", true)

	size, filetype := FileInfo(fileName, link)
	a := strconv.FormatInt(size, 10)
	doLogging("Length:"+a+CalcSize(size)+"["+filetype+"]", true)

	doLogging("Saving to file: "+fileName, true)
	doLogging("", true)

	return resp, a
}

// Write the response of the GET request to file
func writeToFile(directory, fileName string, resp *http.Response) (elapsed time.Duration, data int64) {
	var file *os.File
	if Flags.P_Flag == "" {
		createPath("downloads/" + directory)
		file, _ = os.OpenFile("downloads/"+directory+fileName, os.O_CREATE|os.O_WRONLY, 0777)
	} else {
		createPath(directory)
		file, _ = os.OpenFile(directory+fileName, os.O_CREATE|os.O_WRONLY, 0777)
	}
	createPath("downloads/")
	createPath("downloads/" + directory)

	//var bar progressbar.ProgressBar
	start := time.Now()
	defer file.Close()

	bufferedWriter := bufio.NewWriterSize(file, bufSize)
	if Flags.B_Flag {
		if Flags.RateLimit_Flag != "" {
			limit := ConvertLimit(Flags.RateLimit_Flag)
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
			limit := ConvertLimit(Flags.RateLimit_Flag)
			bucket := ratelimit.NewBucketWithRate(float64(limit), int64(limit))
			data, _ = io.Copy(io.MultiWriter(bufferedWriter, bar), ratelimit.Reader(resp.Body, bucket))
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
		goto exit
	}

	contentType = http.DetectContentType(buf)
exit:
	return size, contentType
}
