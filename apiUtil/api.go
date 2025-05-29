package apiUtil

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var GClient *http.Client

// Initialize the http.Client with a custom Transport
func Init() {

	transport := &http.Transport{
		MaxIdleConns:        20,               // Total number of idle connections
		IdleConnTimeout:     30 * time.Second, // Timeout for idle connections
		MaxIdleConnsPerHost: 5,                // Max idle connections per host
	}

	GClient = &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second, // Timeout for HTTP requests
	}

}

type HeaderDetails struct {
	Key   string
	Value string
}

func Api_call(url string, methodType string, jsonData string, header []HeaderDetails, Source string) (string, error) {
	var body []byte
	var err error
	var request *http.Request
	var response *http.Response
	if methodType != "GET" {
		request, err = http.NewRequest(strings.ToUpper(methodType), url, bytes.NewBuffer([]byte(jsonData)))
	} else {
		request, err = http.NewRequest(strings.ToUpper(methodType), url, nil)
	}

	if err != nil {
		return "", err
	} else {
		if len(header) > 0 {
			for i := 0; i < len(header); i++ {
				request.Header.Set(header[i].Key, header[i].Value)
			}
		}
		// client := &http.Client{}
		response, err = GClient.Do(request)
		if err != nil {
			return "", err
		} else {
			defer response.Body.Close()

			body, err = ioutil.ReadAll(response.Body)
			if err != nil {
				return "", err
			}
		}
	}
	return string(body), nil

}
