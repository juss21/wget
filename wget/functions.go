package wget

import (
	"fmt"
	"net/http"
	"os"
)

// Check if we received an error on our last function call
func errorChecker(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// Make the GET request to a url, return the response
func getResponse(url string) *http.Response {
	fmt.Println("Resolving")
	tr := new(http.Transport)
	fmt.Println("Connecting")
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	errorChecker(err)
	fmt.Println("")
	return resp
}
