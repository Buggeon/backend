package s3storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client   *s3.Client
	endpoint string
	bucket   string
}

func NewS3Storage() *S3Storage {

	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	region := os.Getenv("S3_REGION")

	if endpoint == "" || accessKey == "" || secretKey == "" {
		fmt.Errorf("Missing S3 environment variables")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)

	if err != nil {
		return nil
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	storage := &S3Storage{
		client:   client,
		endpoint: endpoint,
		bucket:   "buggeon",
	}

	storage.ensureBucket(context.Background())

	return storage

}

func (s *S3Storage) ensureBucket(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err == nil {
		return nil
	}

	_, err = s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	return s.putPublicPolicy(ctx)
}

func (s *S3Storage) putPublicPolicy(ctx context.Context) error {
	policy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect":    "Allow",
				"Principal": "*",
				"Action":    "s3:GetObject",
				"Resource":  fmt.Sprintf("arn:aws:s3:::%s/*", s.bucket),
			},
		},
	}
	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	_, err = s.client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(s.bucket),
		Policy: aws.String(string(policyJSON)),
	})
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}
	return nil
}

func (s *S3Storage) Upload(ctx context.Context, key string, body io.Reader, contentType string) (string, error) {

	if key == "" {
		return "", fmt.Errorf("Key cannot be empty")
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: &contentType,
	}

	_, err := s.client.PutObject(ctx, input)

	if err != nil {
		return "", fmt.Errorf("Failed to upload file: %w", err)
	}

	return s.GetUrl(s.bucket, key), nil

}

func (s *S3Storage) GetUrl(bucket, key string) string {

	return s.endpoint + "/" + bucket + "/" + key

}
