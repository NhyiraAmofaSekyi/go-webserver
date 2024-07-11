package aws

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func ListBucketOBJ() error {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String("arn:aws:s3:eu-north-1:049991758581:accesspoint/test2"),
	})

	if err != nil {

		log.Fatal(err)
		return err
	}

	log.Println("first page results:")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
	return nil
}

func GetObject(name string, bucket string) (*s3.GetObjectOutput, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatal(err)
	}

	if name == "" {
		return nil, fmt.Errorf("object name cannot be empty")
	}
	if bucket == "" {
		return nil, fmt.Errorf("bucket name cannot be empty")
	}

	// Create an S3 client from the configuration
	client := s3.NewFromConfig(cfg)

	// Attempt to get the object from the S3 bucket
	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket:       aws.String(bucket),          // Specify the bucket name
		Key:          aws.String(name),            // Specify the object key
		RequestPayer: types.RequestPayerRequester, // Set who pays for the request
	})

	if err != nil {
		// If there is an error, return nil for the object and the error
		return nil, err
	}

	// If no error, return the response object and nil for the error
	return resp, nil
}

func UploadFile(bucketName string, objectKey string, fileName string, contentType string) error {

	if bucketName == "" {
		return fmt.Errorf("bucket name cannot be empty")
	}
	if objectKey == "" {
		return fmt.Errorf("object key cannot be empty")
	}
	if fileName == "" {
		return fmt.Errorf("file name cannot be empty")
	}
	if contentType == "" {
		return fmt.Errorf("content type cannot be empty")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatal(err)
	}

	// Create an S3 client from the configuration
	client := s3.NewFromConfig(cfg)

	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	} else {
		defer file.Close()
		println(contentType)
		_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(objectKey),
			Body:        file,
			ContentType: aws.String(contentType),
		})
		if err != nil {
			log.Printf("Couldn't upload file to %v:%v. Here's why: %v\n",
				bucketName, objectKey, err)
		}
	}
	return err
}

func Upload(bucketName string, objectKey string, file multipart.File, contentType string) error {

	if bucketName == "" {
		return fmt.Errorf("bucket name cannot be empty")
	}
	if objectKey == "" {
		return fmt.Errorf("object key cannot be empty")
	}
	if contentType == "" {
		return fmt.Errorf("content type cannot be empty")
	}
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatal(err)
	}

	// Create an S3 client from the configuration
	client := s3.NewFromConfig(cfg)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Printf("Couldn't upload file to %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return err
}
