package wget

import "net/http"

func mirrorResponse(url string) *http.Response {

	//client
	req, err := http.NewRequest("GET", url, nil)
	errorHandler(err, true)
	resp, err3 := http.DefaultClient.Do(req)
	errorHandler(err3, true)
	return resp
}
