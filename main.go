package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
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

	lastRunTime := getLastRunTime()

	accessToken, err := getAccessToken(clientID, clientSecret)
	if err != nil {
		fmt.Println("Access error")
		return
	}

	assests, err := getUpdatedAssets(accessToken, lastRunTime)
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
	uploader := s3manager.NewUploader(s3.NewFromConfig(cfg))

	for _, asset := range assests {
		err := uploadToS3(uploader, asset)
		if err != nil {
			fmt.Printf("S3 upload error: %v\n", err)
		}
	}
}

func getAccessToken(clientID, clientSecret string) (string, error) {
	reqBody := fmt.Sprintf(`{
        "grant_type": "client_credentials",
        "client_id": "%s",
        "client_secret": "%s"
    }`, clientID, clientSecret)

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

func getLastRunTime() time.Time {
	data, err := ioutil.ReadFile(lastRunFile)
	if err != nil {
		// Varsayılan olarak 24 saat öncesini kullan
		return time.Now().Add(-24 * time.Hour)
	}

	lastRunTime, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return time.Now().Add(-24 * time.Hour)
	}

	return lastRunTime
}

func updateLastRunTime() {
	currentTime := time.Now().Format(time.RFC3339)
	err := ioutil.WriteFile(lastRunFile, []byte(currentTime), 0644)
	if err != nil {
		fmt.Printf("Update error for last run time: %v\n", err)
	}
}

func getUpdatedAssets(accessToken string, since time.Time) ([]Asset, error) {
	query := fmt.Sprintf(`{
        "query": {
            "leftOperand": {
                "property": "modifiedDate",
                "simpleOperator": "greaterThan",
                "value": "%s"
            }
        }
    }`, since.Format(time.RFC3339))

	req, err := http.NewRequest("POST", assestsURL, strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Err req: %s", resp.Status)
	}

	var assets struct {
		Items []Asset `json:"items"`
	}
	err = json.NewDecoder(resp.Body).Decode(&assets)
	if err != nil {
		return nil, err
	}

	return assets.Items, nil
}

func uploadToS3(uploader *s3manager.Uploader, asset Asset) error {
	content := []byte(asset.Content)
	fileName := fmt.Sprintf("%d_%s.json", asset.ID, asset.Name)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(content),
	})

	return err
}
