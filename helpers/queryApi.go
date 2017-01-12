package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

func GetAPI(url string) ([]byte, error) {
	println("Getting from: " + url)
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctxhttp.Do(ctx, nil, req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Error fetching data!")
	}

	return ioutil.ReadAll(resp.Body)
}

func PutAPI(url string, data map[string]interface{}) ([]byte, error) {
	println("Posting to: " + url)
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctxhttp.Do(ctx, nil, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func DeleteAPI(url string) ([]byte, error) {
	println("Deleting: " + url)
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctxhttp.Do(ctx, nil, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
