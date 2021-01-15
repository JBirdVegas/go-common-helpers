package helpers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

func IsSuccessfulHttpStatus(response http.Response) bool {
	return 200 <= response.StatusCode && response.StatusCode < 300
}

func HttpWithHeaders(method string, url string, data []byte, headers *map[string]string) ([]byte, error) {
	request, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for key, value := range *headers {
			request.Header.Add(key, value)
		}
	}
	client := &http.Client{Transport: &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
		DisableCompression:  true,
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
	}}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return all, nil
}
