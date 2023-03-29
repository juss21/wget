package wget

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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
		if Flags.P_Flag != "" {
			givenpath = Folder(Flags.P_Flag)
		}
		if Flags.O_Flag != "" {
			tempFile = Flags.O_Flag
		}

		download_started := time.Now().Format("--2006-01-02 15:04:05--")
		doLogging(download_started+"\t"+url, true)

		if !Flags.Mirror_Flag {
			go startDownload(url, shorturl, tempFile, givenpath, httpmethod)
		} else {
			file, content := startMirroring(url, httpmethod, "", "")
			links, images := GetLinksFromTemp(file, content)
			
			NewLinks := AppendLinks(links, images, url)
			fmt.Println(NewLinks)
			//Flags.Links = NewLinks
			DownloadLinks(NewLinks, httpmethod)
		}
		if i >= 2 {
			break
		}
	}
	wg.Wait()
}



func startDownload(url, shorturl, filename, givenpath, httpMethod string) {

	response, filesize, FileNameFromResp := getResponse(url, httpMethod, shorturl, filename, givenpath)
	filename = FileNameFromResp
	elapsed, data := writeToFile(givenpath, filename, response)

	h, _ := time.ParseDuration(elapsed.String())
	size, _ := strconv.Atoi(filesize)
	var AvgDown float64
	if h.Milliseconds() > 1 && h.Seconds() < 1 {
		AvgDown = float64(size) * (h.Seconds()) / 1000
	} else if h.Seconds() < 1 {
		AvgDown = float64(size) / (h.Seconds()) / 10000000
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
	if !strings.HasSuffix(givenpath, "/") {
		givenpath += "/"
	}
	doLogging(" MB/s"+"\nFile location: '"+givenpath+filename+"'", true)

	wg.Done()
}
