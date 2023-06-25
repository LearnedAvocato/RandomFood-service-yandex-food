package network

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Jeffail/gabs"
)

func DoGetRequest(url string) (*gabs.Container, error) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		return nil, err
	}

	return jsonParsed, err
}

func DoPostRequest(url, payload string) (*gabs.Container, error) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Post(url, "	application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		return nil, err
	}

	return jsonParsed, err
}
