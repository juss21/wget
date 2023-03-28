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

	"github.com/juju/ratelimit"
	"github.com/schollz/progressbar/v3"
)

const (
	bufSize = 1024 * 8
)

var WgetLogNum int

var path = createPath("downloads/")

func Run() {
	if Flags.B_Flag {
		os.Truncate("wget-log", 0)
		fmt.Println("Output will be written to \"wget-log\".")
	}
	// looping all the links saved in Flags
	for l := 0; l < len(Flags.Links); l++ {

		url, shorturl, givenfilename, givenpath, httpmethod := sliceUrl(Flags.Links[l])
		if Flags.P_Flag != "" {
			givenpath = Folder(Flags.P_Flag)
		}
		download_started := time.Now().Format("--2006-01-02 15:04:05--")
		doLogging(download_started, false)
		doLogging("\t"+url, true)
		tempFile := givenfilename
		if Flags.O_Flag != "" {
			tempFile = Flags.O_Flag
		}
		response, FileSizeString := getResponse(url, httpmethod, shorturl, tempFile)
		//size, _ := FileInfo(tempFile, url)
		elapsed, DownloadedData := writeToFile(givenpath+"/", tempFile, response)
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
		doLogging(" MB/s"+"\nFile location: '"+givenpath+tempFile+"'", true)

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

		f, _ := os.OpenFile("wget-log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

		defer f.Close()

		if _, err := f.WriteString(input); err != nil {
			panic(err)
		}
		if newline {
			f.WriteString("\n")
		}
		WgetLogNum++
	}
}
func FileExists(filepath string) bool {

	fileinfo, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false
	}
	// Return false if the fileinfo says the file path is a directory.
	return !fileinfo.IsDir()
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
func Folder(path string) string {
	if strings.HasPrefix(path, "~") {
		path = strings.TrimPrefix(path, "~")
		home := os.Getenv("HOME")
		path = home + path
		err := os.MkdirAll(path, 0755)
		errorChecker(err)
	} else {
		path = "downloads/" + path
		err := os.MkdirAll(path, 0755)
		errorChecker(err)
	}
	return path
}

func ConvertLimit(base string) int {
	var limit int
	if strings.HasSuffix(strings.ToLower(base), "b") {
		limit, _ = strconv.Atoi(strings.TrimSuffix(base, "b"))
	} else if strings.HasSuffix(strings.ToLower(base), "k") {
		limit, _ = strconv.Atoi(strings.TrimSuffix(base, "k"))
		limit *= 1024
	} else if strings.HasSuffix(strings.ToLower(base), "m") {
		limit, _ = strconv.Atoi(strings.TrimSuffix(base, "m"))
		limit *= 1048576
	}
	return limit
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

	//var bar progressbar.ProgressBar
	start := time.Now()
	defer file.Close()

	bufferedWriter := bufio.NewWriterSize(file, bufSize)


	if Flags.B_Flag {
		if Flags.RateLimit_Flag != "" {
			limit:= ConvertLimit(Flags.RateLimit_Flag)
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
			limit:= ConvertLimit(Flags.RateLimit_Flag)
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
