package wget

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func errorHandler(err error, fatal bool) {
	if err != nil {
		fmt.Println(err)
		if fatal {
			os.Exit(0)
		}
	}
}

func Folder(path string) string {
	if strings.HasPrefix(path, "~") {
		path = strings.TrimPrefix(path, "~")
		home := os.Getenv("HOME")
		path = home + path
		err := os.MkdirAll(path, 0755)
		errorHandler(err, true)
	} else {
		path = "downloads/" + path
		err := os.MkdirAll(path, 0755)
		errorHandler(err, true)
	}
	return path
}

// function for logging
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
	}
}

// shorten time duration
func shortTimeDur(d time.Duration) string {
	s := d.String()
	v, err := strconv.ParseFloat(s[:len(s)-2], 64)
	if err != nil {
		fmt.Println(err)
	}
	return fmt.Sprintf("%.2f", v) + s[len(s)-2:]
}

// create path if not exists
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

// get port from httpmethod
func GetPortFromHttpMethod(s string) (port string) {
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

// calc byte value
func CalcSize(numb int64) string {
	var output string

	numstr := strconv.FormatInt(numb, 10)
	num, err := strconv.Atoi(numstr)
	errorHandler(err, true)
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

// Read data from url
func sliceUrl(url string) (rurl, cleanurl, givenfilename, givenpath, httpmethod string) {
	var rebuilt []string

	split := strings.Split(url, "/")

	for i := 0; i < len(split); i++ {
		if split[i] != "" {
			rebuilt = append(rebuilt, split[i])
		}
	}

	rurl = url

	if len(split) > 4 {
		if split[3] != "" {

			givenpath = "downloads/" + rebuilt[2]
			for j := len(split) - 1; j > 2; j-- {
				if strings.Contains(split[j], ".") {
					givenfilename = split[j]
					break
				}
			}
		}
	} else {
		givenpath = "downloads/"
		givenfilename = split[len(split)-1]
	}
	cleanurl = rebuilt[1]
	httpmethod = rebuilt[0]

	return rurl, cleanurl, givenfilename, givenpath, httpmethod
}

func AppendLinks(links, images []string, baseurl string) (FinalLinks []string) {
	for _, link := range links {
		if !strings.HasPrefix(link, "http") {
			FinalLinks = append(FinalLinks, baseurl+link)
		}
	}
	for _, img := range images {
		if !strings.HasPrefix(img, "http") {
			FinalLinks = append(FinalLinks, baseurl+img)
		} else {
			FinalLinks = append(FinalLinks, img)
		}
	}

	for fl := 0; fl < len(FinalLinks); fl++ {
		FinalLinks[fl] = strings.Replace(FinalLinks[fl], "/./", "/", -1)
	}

	return
}
