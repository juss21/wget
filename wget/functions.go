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
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
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
	cleanurl = rebuilt[1]
	if len(split) > 4 {
		givenfilename = rebuilt[2]

		if len(split) != 4 {
			givenpath = "downloads/" + rebuilt[2]
			givenfilename = rebuilt[3]
		}
	} else {
		givenpath = "downloads/"
	}
	httpmethod = rebuilt[0]

	return rurl, cleanurl, givenfilename, givenpath, httpmethod
}
