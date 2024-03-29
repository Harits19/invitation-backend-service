package bucket

import (
	"context"
	"fmt"
	"main/common/constan"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var client *s3.S3

var s3Config *aws.Config

type Bucket struct {
	bucketName string
}

func New(name string) Bucket {

	return Bucket{
		bucketName: name,
	}
}

func InitConnection() {
	s3Config = &aws.Config{
		Region:      aws.String(constan.ENV.OBJECT_STORAGE_REGION),
		Endpoint:    aws.String(constan.ENV.OBJECT_STORAGE_ENDPOINT),
		Credentials: credentials.NewStaticCredentials(constan.ENV.OBJECT_STORAGE_ACCESS_KEY, constan.ENV.OBJECT_STORAGE_SECRET_KEY, ""),
	}
	fmt.Println("start init storage object")
	resultSession := session.Must(session.NewSession(s3Config))

	client = s3.New(resultSession)
}

func (bucket Bucket) UploadFile(path string, fileHeader multipart.FileHeader) (*string, error) {
	BucketName := bucket.bucketName
	timeout := 10 * time.Minute

	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	if cancelFn != nil {
		defer cancelFn()
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	_, err = client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:             aws.String(BucketName),
		Key:                aws.String(path),
		Body:               file,
		ACL:                aws.String("public-read"),
		ContentType:        aws.String(""),
		ContentDisposition: aws.String("inline"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
		}
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/%s", *s3Config.Endpoint, BucketName, path)
	fmt.Println("successfully uploaded file to ", url)
	return &url, nil
}

func (bucket Bucket) CreateBucket() error {
	bucketName := bucket.bucketName
	// BucketName := bucket.BucketName
	result, err := client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		fmt.Println("error when create bucket", err)
		return err
	}
	fmt.Println("success create bucket on location ", *result.Location)
	return nil
}

func (bucket Bucket) getPublicUrl() {
	BucketName := bucket.bucketName

	res, _ := client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(BucketName),
	})
	for _, object := range res.Contents {
		fmt.Printf("Found object: %s, size: %d\n", *object.Key, *object.Size)
		req, _ := client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(BucketName),
			Key:    aws.String(*object.Key),
		})

		urlStr, _ := req.Presign(15 * time.Minute)

		fmt.Println("urlStr", urlStr)
	}
}

func (bucket Bucket) deleteAllBucketObject() {
	BucketName := bucket.bucketName

	res, _ := client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(BucketName),
	})
	for _, object := range res.Contents {
		fmt.Printf("Found object: %s, size: %d\n", *object.Key, *object.Size)

		_, _ = client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(BucketName),
			Key:    aws.String(*object.Key),
		})
	}
}

func (bucket Bucket) GetListOfBucketFile() (*s3.ListObjectsV2Output, error) {
	BucketName := bucket.bucketName

	return client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(BucketName),
	})

}

func (bucket Bucket) FindObjectByName(name string) {
	client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket.bucketName),
		Key:    aws.String(name),
	})
}

func (bucket Bucket) DeleteObjectByName(name string) error {

	fmt.Println("start delete file with key", name)
	_, err := client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket.bucketName),
		Key:    aws.String(name),
	})

	return err
}

func generateFileName(prefix string, index int) string {
	return fmt.Sprintf("%s/%s.%d", "assets", prefix, index)
}

func (bucket Bucket) SaveToStorage(file *multipart.FileHeader, prefix string, index int) (string, error) {

	fileName := strings.Split(file.Filename, ".")
	fileExtension := fileName[len(fileName)-1]

	newFileName := generateFileName(prefix, index)

	uniqueId := time.Now().Unix()

	filePath := fmt.Sprintf("%s_%d.%s", newFileName, uniqueId, fileExtension)

	url, err := bucket.UploadFile(filePath, *file)

	return *url, err
}

func (bucket Bucket) FindAndDelete(files *s3.ListObjectsV2Output, prefix string, index int) error {
	prefix = generateFileName(prefix, index)
	for _, file := range files.Contents {
		if strings.Contains(*file.Key, prefix) {

			fmt.Println("detected file with prefix", prefix)
			fmt.Println("file name", *file.Key)

			err := bucket.DeleteObjectByName(*file.Key)

			return err

		}
	}
	return nil
}
