package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-openapi/jsonpointer"
	"github.com/pkg/errors"
)

func getDSN(params Parameters) (string, error) {
	switch params.Source {
	case "", "env":
		dsn := os.Getenv("AIRPLANE_DSN")
		if dsn != "" {
			return dsn, nil
		} else {
			return "", errors.New("Expected an AIRPLANE_DSN environment variable to be set.")
		}
	case "s3":
		return getDSNFromS3(params)
	default:
		return "", fmt.Errorf("Unknown --source: %s", params.Source)
	}
}

func getDSNFromS3(params Parameters) (string, error) {
	sess, err := session.NewSession()
	if err != nil {
		return "", errors.Wrap(err, "starting session")
	}

	// Parse the bucket and path from the S3 url:
	trimmed := strings.TrimPrefix(params.S3URL, "s3://")
	components := strings.SplitN(trimmed, "/", 2)
	if len(components) != 2 || components[0] == "" || components[1] == "" {
		return "", errors.Errorf("invalid s3 url: expected s3://bucket-name/file-path got %s", params.S3URL)
	}
	bucket, key := components[0], components[1]

	buffer := aws.NewWriteAtBuffer([]byte{})
	if _, err = s3manager.NewDownloader(sess).Download(buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return "", errors.Wrap(err, "downloading s3 object")
	}

	var contents interface{}
	if err = json.Unmarshal(buffer.Bytes(), &contents); err != nil {
		return "", errors.Wrap(err, "s3 file contains invalid JSON")
	}

	ptr, err := jsonpointer.New(params.S3JSONPointer)
	if err != nil {
		return "", errors.Wrap(err, "invalid JSON pointer")
	}

	dsn, typ, err := ptr.Get(contents)
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve JSON pointer")
	}
	if typ != reflect.String {
		return "", errors.Wrapf(err, "JSON pointer resolved to a non-string value: got %s", typ)
	}

	if dsn == "" {
		return "", errors.New("s3 file contains empty DSN")
	}

	return dsn.(string), nil
}
