package components

import (
	"context"
	"io/ioutil"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

func QueryAPI(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	go func() {
		cancel()
	}()

	resp, err := ctxhttp.Get(ctx, nil, url)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}
