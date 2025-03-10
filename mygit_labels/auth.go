package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	grant_type    = "urn:ietf:params:oauth:grant-type:device_code"
	auth_pending  = "authorization_pending"
	slow_down     = "slow_down"
	expired_token = "token_expired"
	access_denied = "access_denied"
	access_token  = "access_token"
	refresh_token = "refresh_token"
)

type loginEssentials struct {
	Device_code      string        `json:"device_code"`
	User_code        string        `json:"user_code"`
	Verification_uri string        `json:"verification_uri"`
	Interval         time.Duration `json:"interval"`
	// Expires_in       time.Time     `json:"expires_in,omitempty"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Error        string `json:"error"`
}

type User struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Url   string `json:"html_url"`
	// ViewType string `json:"user_view_type,omitempty"`
}

func makeRequest(req *http.Request) (int, []byte, error) {
	// set default headers
	req.Header.Set("User-Agent", "github_api_go_client")
	req.Header.Set("X-GitHup-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	status := resp.StatusCode

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	log.Printf("(%v) \n", resp.Request.URL.Path)
	return status, data, nil
}

func GetUsername(token string) (string, error) {
	url := BASE_API_URL + "user"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	fmt.Printf("Authorization: %s\n", req.Header.Get("Authorization"))
	_, body, err := makeRequest(req)
	if err != nil {
		return "", err
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	var name string
	switch resp["name"].(type) {
	case string:
		name = resp["name"].(string)
	}
	return name, nil
}

func Login(client_id string) (string, string, error) {
	// get_request_device
	code_url := BASE_AUTH_URL + "login/device/code"
	ctx, canc := context.WithTimeout(context.Background(), 60*time.Second)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, code_url, nil)
	defer ctx.Done()
	defer canc()
	if err != nil {
		return "", "", err
	}
	q := req.URL.Query()
	q.Set("client_id", client_id)
	req.URL.RawQuery = q.Encode()

	_, body, err := makeRequest(req)
	if err != nil {
		return "", "", err
	}
	var params loginEssentials
	if err := json.Unmarshal(body, &params); err != nil {
		return "", "", err
	}

	fmt.Printf("Please visit %s,\n and enter your code (%s).\n", params.Verification_uri, params.User_code)

	login_url := BASE_AUTH_URL + "login/oauth/access_token"
	resp, err := pullRequestToken(login_url, client_id, params)
	if err != nil {
		return "", "", err
	}
	return resp.AccessToken, resp.RefreshToken, nil
}

func Refresh(client_id, refreshToken string) (string, string, error) {
	var data LoginResponse
	login_url := BASE_AUTH_URL + "login/oauth/access_token"
	req, err := http.NewRequest(http.MethodPost, login_url, nil)
	if err != nil {
		return "", "", err
	}
	q := req.URL.Query()
	q.Set("client_id", client_id)
	q.Set("grant_type", refresh_token)
	q.Set("refresh_token", refreshToken)
	req.URL.RawQuery = q.Encode()

	_, resp, err := makeRequest(req)
	if err := json.Unmarshal(resp, &data); err != nil {
		return "", "", err
	}
	log.Printf("received %+v", data)
	return data.AccessToken, data.RefreshToken, nil
}

func pullRequestToken(login_url, client_id string, params loginEssentials) (LoginResponse, error) {
	var resp LoginResponse
	req, err := http.NewRequest(http.MethodPost, login_url, nil)
	if err != nil {
		return resp, err
	}
	q := req.URL.Query()
	q.Set("client_id", client_id)
	q.Set("device_code", params.Device_code)
	q.Set("grant_type", grant_type)
	req.URL.RawQuery = q.Encode()

	// wait for pull a login request
	for {
		time.Sleep(params.Interval * time.Second)
		_, body, err := makeRequest(req)
		if err != nil {
			return resp, err
		}
		var resp LoginResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return resp, err
		}
		errRes := resp.Error
		switch errRes {
		case auth_pending:
			params.Interval++
			continue
		case slow_down:
			params.Interval += 5
			continue
		case expired_token:
		case access_denied:
			return resp, errors.New(errRes)
		}
		if errRes != "" {
			log.Fatalf("App is not working properly.\nPlease try again!\n%s", errRes)
		}
		if resp.AccessToken != "" {
			log.Println(resp)
			return resp, nil
		}
	}
}
