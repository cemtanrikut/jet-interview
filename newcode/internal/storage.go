package newcode

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
)

// UploadAssetsToS3 uploads assets to the specified S3 bucket.
func UploadAssetsToS3(bucketName string, assets []Asset) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("AWS upload err: %v", err)
	}

	uploader := s3manager.NewUploader(s3.NewFromConfig(cfg))

	for _, asset := range assets {
		content := []byte(asset.Content)
		fileName := fmt.Sprintf("%d_%s.json", asset.ID, asset.Name)

		_, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fileName),
			Body:   bytes.NewReader(content),
		})
		if err != nil {
			return fmt.Errorf("Asset %s doens't upload: %v", fileName, err)
		}
	}

	return nil
}
