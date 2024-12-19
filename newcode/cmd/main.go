package main

import (
	"time"

	newcode "jet-interview/newcode"

	logrus "github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	cfg := newcode.LoadConfig()

	// Salesforce access token
	token, err := newcode.internal.GetAccessToken(cfg.ClientID, cfg.ClientSecret)
	if err != nil {
		log.Fatalf("Access token problem: %v", err)
	}

	// Get last run time
	lastRunTime := newcode.internal.GetLastRunTime()

	// Fetch updated assets
	assets, err := newcode.internal.FetchUpdatedAssets(token, lastRunTime)
	if err != nil {
		log.Fatalf("Fetching problem: %v", err)
	}

	// upload items to AWS S3
	err = newcode.internal.UploadAssetsToS3(cfg.BucketName, assets)
	if err != nil {
		log.Fatalf("S3 uploding problem : %v", err)
	}

	// Son çalıştırma zamanını güncelle
	newcode.internal.UpdateLastRunTime(time.Now())

	log.Info("Done!")
}
