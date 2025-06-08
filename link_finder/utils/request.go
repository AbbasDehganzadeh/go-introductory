package utils

import (
	"io"
	"net/http"
)

type Status struct {
	Code   int
	Reason string
}

func MakeRequest(url string) (Status, []byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Status{Reason: "Unknown Url"}, nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Status{Reason: "Unknown Network"}, nil, err
	}
	defer resp.Body.Close()
	status := getStatus(resp)
	body, err := getBody(resp)
	if err != nil {
		return Status{Reason: "Unknown Body"}, nil, err
	}
	return status, body, err
}

func getStatus(resp *http.Response) Status {
	return Status{
		Code:   resp.StatusCode,
		Reason: resp.Status,
	}
}

func getBody(resp *http.Response) ([]byte, error) {
	//? use readall!
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
