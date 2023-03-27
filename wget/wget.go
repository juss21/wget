package wget

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
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
	fmt.Println("\t\t" + url)

	response := getResponse(url, url_split)
	writeToFile(url_filename, response, url_split)

	// fmt.Println(url_split, url_httpstatus, url_webaddress, url_directory, url_filename)

	// DownloadFile(url_directory+"/"+url_filename, url)

}

// Write the response of the GET request to file
func writeToFile(fileName string, resp *http.Response, url_split []string) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0777)
	errorChecker(err)
	defer file.Close()
	bufferedWriter := bufio.NewWriterSize(file, bufSize)
	errorChecker(err)
	_, err = io.Copy(bufferedWriter, resp.Body)
	errorChecker(err)
}
