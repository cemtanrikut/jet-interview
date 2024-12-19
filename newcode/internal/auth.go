package newcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const tokenURL = "https://YOUR_SUBDOMAIN.auth.marketingcloudapis.com/v2/token"

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetAccessToken retrieves an access token from Salesforce.
func GetAccessToken(clientID, clientSecret string) (string, error) {
	reqBody := fmt.Sprintf(`{
        "grant_type": "client_credentials",
        "client_id": "%s",
        "client_secret": "%s"
    }`, clientID, clientSecret)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Req err: %s", resp.Status)
	}

	var tokenResp AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}
