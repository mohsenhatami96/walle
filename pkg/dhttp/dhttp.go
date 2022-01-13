package dhttp

import (
	"io/ioutil"
	"net/http"
)

func Getter(url string, token *string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return make([]byte, 0), err
	}
	if token != nil {
		req.Header.Set("PRIVATE-TOKEN", *token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return make([]byte, 0), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make([]byte, 0), err
	}

	return body, err
}
