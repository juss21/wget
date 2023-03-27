package wget

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	bufSize = 1024 * 8
)

func Run(url, filename string) {
	url_split := strings.Split(url, "/")
	// url_httpstatus := url_split[0]
	// url_webaddress := url_split[2]
	//	url_directory := ""
	url_filename := url_split[3]
	if len(url_split) == 5 {
		//		url_directory = url_split[3]
		url_filename = url_split[4]
	}

	download_started := time.Now().Format("--2006-01-02 15:04:05--")
	fmt.Print(download_started)
	fmt.Println("\t" + url)

	response := getResponse(url, url_split)
	size, _ := FileInfo(url_split[4], url)
	elapsed := writeToFile(url_filename, response, url_split)
	h, _ := time.ParseDuration(elapsed.String())
	var AvgDown float64
	if h.Seconds() < 1 {
		AvgDown = float64(size) * (h.Seconds()) / 10
	} else {
		AvgDown = float64(size) / (h.Seconds()) / 1000000
	}
	//fmt.Println(elapsed, size, h.Seconds())

	fmt.Println(time.Now().Format("2006-01-02 15:04:05 - Download completed! "), "Time elapsed:", elapsed,
	 "Average download speed:", math.Round(AvgDown*10)/10, "MB/s")

	// fmt.Println(url_split, url_httpstatus, url_webaddress, url_directory, url_filename)

	// DownloadFile(url_directory+"/"+url_filename, url)

}

// Write the response of the GET request to file
func writeToFile(fileName string, resp *http.Response, url_split []string) (elapsed time.Duration) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0777)
	errorChecker(err)

	start := time.Now()
	defer file.Close()

	bufferedWriter := bufio.NewWriterSize(file, bufSize)
	errorChecker(err)
	bar := progressbar.DefaultBytes( 
		resp.ContentLength,
		"Downloading: "+url_split[4],
	)
	_, err = io.Copy(io.MultiWriter(bufferedWriter, bar), resp.Body)
	errorChecker(err)
	t := time.Now()
	elapsed = t.Sub(start)

	return elapsed
}
