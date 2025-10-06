package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"time"
)

type bucket string

func (b bucket) string() string {
	return string(b)
}

func (s *service) getPresignURL(presignClient *s3.PresignClient, bucket bucket, object string, expired time.Duration) (string, error) {
	presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket.string()),
		Key:    aws.String(object),
	}, func(options *s3.PresignOptions) {
		options.Expires = expired
	})
	if err != nil {
		return "", err
	}
	return presignResult.URL, err
}

func (s *service) UploadFile(client *s3.Client, bucket bucket, file io.Reader, path string) error {
	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket.string()),
		Key:    aws.String(path),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to %s/%s: %w", bucket, path, err)
	}
	return nil
}

func (s *service) CopyFileToAnotherBucket(client *s3.Client, oldBucket bucket, newBucket bucket, path string) error {
	_, err := client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(newBucket.string()),
		CopySource: aws.String(fmt.Sprintf("%s/%s", oldBucket.string(), path)),
		Key:        aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to copy object for soft delete: %w", err)
	}
	return nil
}

func (s *service) DeleteFile(client *s3.Client, bucket bucket, path string) error {
	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket.string()),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %s/%s: %w", bucket, path, err)
	}
	return nil
}

func (s *service) ReadFile(client *s3.Client, bucket bucket, path string) ([]byte, error) {
	output, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket.string()),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object %s/%s: %w", bucket, path, err)
	}
	defer output.Body.Close()

	content, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %w", err)
	}
	return content, nil
}

func (s *service) ListFiles(client *s3.Client, bucket bucket) error {
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:                   aws.String(bucket.string()),
		ContinuationToken:        nil,
		Delimiter:                nil,
		EncodingType:             "",
		MaxKeys:                  nil,
		OptionalObjectAttributes: nil,
		Prefix:                   nil,
		RequestPayer:             "",
		StartAfter:               nil,
	})
	if err != nil {
		return fmt.Errorf("failed to list objects in bucket %s: %w", bucket, err)
	}

	fmt.Println(output)

	return nil
}
