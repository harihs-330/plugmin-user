package apikit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func APIRequest(
	ctx context.Context,
	url, method string,
	headers, body map[string]interface{},
	result interface{}) (int, error) {

	var reader io.Reader
	if body != nil {
		bs, err := json.Marshal(body)
		if err != nil {
			return 0, err
		}

		reader = bytes.NewReader(bs)
	}

	request, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		log.Println("Unable to connect to the API server", url, err)
		return 0, err
	}

	request.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Add(k, fmt.Sprintf("%s", v))
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	return resp.StatusCode, json.Unmarshal(responseBody, result)
}
