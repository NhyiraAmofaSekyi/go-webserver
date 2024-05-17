package aws

import (
	"context"
	"log"

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
	// Load the AWS default configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		// Return nil for the object and the error
		return nil, err
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
