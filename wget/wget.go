package wget

import (
	"bufio"
	"fmt"
	"math"
	"strconv"
	"strings"

	//"go/format"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	bufSize = 1024 * 8
)

var path = createPath("downloads/")

func Run() {
	// looping all the links saved in Flags
	for l := 0; l < len(Flags.Links); l++ {
		url, shorturl, givenfilename, givenpath, httpmethod := sliceUrl(Flags.Links[l])

		download_started := time.Now().Format("--2006-01-02 15:04:05--")
		if Flags.B_Flag {
			fmt.Println("Output will be written to \"wget-log\".")
			f, err := os.Create("wget-log")
			errorChecker(err)
			defer f.Close()
		}
		doLogging(download_started, false)
		doLogging("\t"+url, true)
		tempFile := givenfilename
		if Flags.O_Flag != "" {
			tempFile = Flags.O_Flag
		}
		response, FileSizeString := getResponse(url, httpmethod, shorturl, tempFile)
		//size, _ := FileInfo(tempFile, url)
		elapsed, DownloadedData := writeToFile(givenpath, tempFile, response)
		h, _ := time.ParseDuration(elapsed.String())
		size, _ := strconv.Atoi(FileSizeString)
		var AvgDown float64
		if h.Seconds() < 1 {
			AvgDown = float64(size) * (h.Seconds()) / 10
		} else {
			AvgDown = float64(size) / (h.Seconds()) / 1000000
		}
		DownloadedDataInt := strconv.FormatInt(DownloadedData, 10)

		doLogging(time.Now().Format("2006-01-02 15:04:05 - Download completed! "), false)
		doLogging("["+string(DownloadedDataInt)+"/"+FileSizeString+"]", true)
		doLogging("Time elapsed: ", false)
		doLogging(shortDur(elapsed), false)
		doLogging(" Average download speed: ", false)
		math := math.Round(AvgDown*10) / 10
		conv := strconv.FormatFloat(math, 'f', -1, 64)
		doLogging(conv, false)
		doLogging(" MB/s"+"\nFile: '"+path+tempFile+"'", true)

	}

}
func shortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}
func doLogging(input string, newline bool) {
	if !Flags.B_Flag {
		if newline {
			fmt.Println(input)
		} else {
			fmt.Print(input)
		}
	} else {
		f, err := os.OpenFile("wget-log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString(input); err != nil {
			panic(err)
		}
		if newline {
			f.WriteString("\n")
		}

	}
}

// check if path exists, if not create
func createPath(path string) string {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return path
		}
	}
	return path
}

// Write the response of the GET request to file
func writeToFile(directory, fileName string, resp *http.Response) (elapsed time.Duration, data int64) {
	createPath("downloads/" + directory)
	file, err := os.OpenFile("downloads/"+directory+"/"+fileName, os.O_CREATE|os.O_WRONLY, 0777)
	errorChecker(err)
	//var bar progressbar.ProgressBar
	start := time.Now()
	defer file.Close()

	bufferedWriter := bufio.NewWriterSize(file, bufSize)
	errorChecker(err)
	if Flags.B_Flag {
		data, err = io.Copy(bufferedWriter, resp.Body)
		errorChecker(err)
	} else {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"Downloading: "+fileName,
		)
		data, err = io.Copy(io.MultiWriter(bufferedWriter, bar), resp.Body)
		errorChecker(err)
	}
	t := time.Now()
	elapsed = t.Sub(start)

	return elapsed, data
}
