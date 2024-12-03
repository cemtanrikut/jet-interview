package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	tokenURL    = ""
	assestsURL  = ""
	bucketName  = ""
	lastRunFile = ""
)

type AccessTokenResponse struct {
	AccessToken string
	ExpiresIn   int
}

type Asset struct {
	ID           int
	Name         string
	Content      string
	ModifiedDate time.Time
}

func main() {
	clientID := os.Getenv("SF_CLIENT_ID")
	clientSecret := os.Getenv("SF_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		fmt.Println("Client ID and Client Secret are empty")
		return
	}

	accessToken, err := getAccessToken(clientID, clientSecret)
	if err != nil {
		fmt.Println("Access error")
		return
	}

	assests, err := getUpdatedAssets()
	if err != nil {
		fmt.Println("Get updated assests error")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("failed to load configuration, %v", err)
		return
	}
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)
}

func getAccessToken(clientID, clientSecret string) (string, error) {
	reqBody := fmt.Sprintf("")

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(reqBody))
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
		return "", fmt.Errorf("HTTP status error. Not OK")
	}

	var tokenResp AccessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil

}

func getUpdatedAssets() (string, error) {
	return "", nil
}
