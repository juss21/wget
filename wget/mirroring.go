package wget

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func mirrorResponse(filename, url string) (r *http.Response, filetype string) {
	_, filetype = FileInfo(filename, url)
	//client
	req, err := http.NewRequest("GET", url, nil)
	errorHandler(err, true)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, err3 := http.DefaultClient.Do(req)
	errorHandler(err3, true)
	return resp, filetype
}
func GetLinksFromTemp(file *os.File, content []byte) (links, images []string) {

	var corndogRegExp, imageRegExp, linksRegExp *regexp.Regexp

	if Flags.X_Flag != "" || Flags.Reject_Flag != "" {
		if Flags.X_Flag == "" {
			Flags.X_Flag = Flags.Reject_Flag
		}
		corndogRegExp = regexp.MustCompile(`url([(].)([^)])([^']*)`)
		imageRegExp = regexp.MustCompile(`src=["']([^"']+)([^` + Flags.X_Flag + `])["']`)
		linksRegExp = regexp.MustCompile(`href=["']([^"']+)([^` + Flags.X_Flag + `])["']`)
	} else {
		corndogRegExp = regexp.MustCompile(`url([(].)([^)])([^']*)`)
		imageRegExp = regexp.MustCompile(`src=["']([^"']+)["']`)
		linksRegExp = regexp.MustCompile(`href=["']([^"']+)["']`)
	}
	//var links, images []string
	corn := corndogRegExp.FindAllStringSubmatch(string(content), -1)
	for _, item := range corn {
		images = append(images, item[3])
	}

	sub := linksRegExp.FindAllStringSubmatch(string(content), -1)
	for _, item := range sub {
		if Flags.X_Flag != "" || Flags.Reject_Flag != "" {
			links = append(links, item[1]+item[2])
		} else {
			links = append(links, item[1])
		}
	}
	subMatchSlice := imageRegExp.FindAllStringSubmatch(string(content), -1)
	for _, item := range subMatchSlice {
		if Flags.X_Flag != "" || Flags.Reject_Flag != "" {
			images = append(images, item[1]+item[2])
		} else {
			images = append(images, item[1])
		}
	}
	return
}

func DownloadLinks(Links []string, url, httpmethod string) {
	fmt.Println(Links)
	for i, url := range Links {
		strings.Replace(Links[i], "/./", "/", -1)
		if Links[i][len(Links[i])-1] == '/' {
			fmt.Println("err", Links[i], url)
			continue
		}
		fmt.Println(i, Links[i])
		wg.Add(1)
		var path, filename string
		split_url := strings.Split(url, "/")
		for j := 3; j < len(split_url); j++ {
			if j == len(split_url)-1 {
				filename += split_url[j]
				break
			}
			filename += split_url[j] + "/"

		}
		if Flags.P_Flag != "" {
			path = Folder(Flags.P_Flag) + split_url[2]
		} else {
			path = "downloads/" + split_url[2]
		}
		startMirroring(url, httpmethod, filename, path)
		if i >= len(Links) {
			break
		}
	}
	wg.Wait()
}

func startMirroring(url, httpmethod, filename, path string) (*os.File, []byte) {

	surl := strings.Split(url, "/")
	FilenameSlice := strings.Split(filename, "/")

	response, file_type := mirrorResponse(filename, url)

	createPath("downloads/")
	createPath("downloads/" + surl[2])
	if len(FilenameSlice) > 1 {
		createPath("downloads/" + surl[2] + "/" + FilenameSlice[0] + "/")
		for i := 0; i < len(FilenameSlice); i++ {
			createPath("downloads/" + surl[2] + "/" + strings.Join(FilenameSlice[:i], "/"))
		}
	}

	var file *os.File
	var erro error
	if filename == "" {
		file, erro = os.OpenFile("downloads/"+surl[2]+"/index.html", os.O_CREATE|os.O_WRONLY, 0o644)
		errorHandler(erro, true)
	} else {
		if strings.Contains(file_type, "text/html") {
			filename += ".html"
		}
		file, erro = os.OpenFile(path+"/"+filename, os.O_CREATE|os.O_WRONLY, 0o644)
		errorHandler(erro, true)
	}

	CopyThat, err1 := ioutil.ReadAll(response.Body)
	errorHandler(err1, true)

	CopyThat = []byte(strings.ReplaceAll(string(CopyThat), "url('/", "url('./"))
	CopyThat = []byte(strings.ReplaceAll(string(CopyThat), `="/`, `="./`))

	bar := progressbar.DefaultBytes(
		bytes.NewReader(CopyThat).Size(),
		"Downloading: "+url,
	)
	w := io.MultiWriter(file, bar)
	w.Write(CopyThat)
	response.Body.Close()
	fmt.Println()

	wg.Done()
	return file, CopyThat
}
