package refresher

import (
	"auth/pkg/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func RefreshToken() (*model.Token, error) {
	for {
		time.Sleep(10 * time.Minute)

		urlAuth := "http://localhost:8081/auth"
		req, err := http.NewRequest("GET", urlAuth, nil)
		if err != nil {
			fmt.Printf("%v", err)
			return nil, fmt.Errorf("failed to create request to Olist ERP API: %w", err)
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("%v", err)
			return nil, fmt.Errorf("failed to make request to Olist ERP API: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("%v", err)
			return nil, fmt.Errorf("returned non-200 status from Tiny ERP API: %s", resp.Status)
		}
		// Parse the response
		var token model.Token
		if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
			fmt.Printf("%v", err)
			return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
		}
		// check if LastUpdate more than 3 hours ago
		parsedTime, err := time.Parse(time.RFC3339, token.Lastupdate)
		if err != nil {
			fmt.Printf("%v", err)
			return nil, fmt.Errorf("failed to parse Lastupdate: %w", err)
		}
		if time.Since(parsedTime) > 180*time.Minute {
			// perform request for a new access token, using refresh token
			token, err := refreshAccessToken(token.RefreshToken)
			if err != nil {
				return nil, fmt.Errorf("failed to refresh access token: %w", err)
			}
			u, err := url.Parse(urlAuth)
			if err != nil {
				fmt.Printf("%v", err)
				return nil, fmt.Errorf("failed to parse URL: %w", err)
			}
			q := u.Query()
			q.Set("key", token.Key)
			q.Set("refresh_token", token.RefreshToken)
			u.RawQuery = q.Encode()

			req, err := http.NewRequest("PUT", u.String(), nil)
			if err != nil {
				fmt.Printf("%v", err)
				return nil, fmt.Errorf("failed to create request to Olist ERP API: %w", err)
			}

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("%v", err)
				return nil, fmt.Errorf("failed to make request to Olist ERP API: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("returned non-200 status from Tiny ERP API: %s", resp.Status)
			}

			fmt.Printf("Token refreshed successfully at %v \n", token.Lastupdate)
		}

	}
}
func refreshAccessToken(refreshToken string) (*model.Token, error) {
	refreshUrl := "https://accounts.tiny.com.br/realms/tiny/protocol/openid-connect/token" // TODO replace by environment variable
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", os.Getenv("CLIENT_ID"))         // Replace with actual client ID
	form.Set("client_secret", os.Getenv("CLIENT_SECRET")) // Replace with actual client secret
	form.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", refreshUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request to Olist ERP API: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Olist ERP API: %w", err)
	}
	defer resp.Body.Close()

	// print response body to check if ok
	// body, err := io.ReadAll(resp.Body)
	// fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("returned non-200 status from Tiny ERP API: %s", resp.Status)
	}
	// Parse the response
	var response model.RefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("%v", err)
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	var token model.Token
	token.Key = response.AccessToken
	token.RefreshToken = response.RefreshToken
	token.Lastupdate = time.Now().Format(time.RFC3339)

	return &token, nil
}
