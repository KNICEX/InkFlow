package web

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type S3Handler struct {
	bucket string
	client *s3.Client
}

func (s S3Handler) RegisterRoutes(server *gin.Engine) {
	//TODO implement me
	panic("implement me")
}

func NewS3Handler(endpoints, region, accessKey, secretKey, bucket string) *S3Handler {
	client := s3.NewFromConfig(aws.Config{
		Region:       region,
		Credentials:  credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		BaseEndpoint: aws.String(endpoints),
	})
	return &S3Handler{
		bucket: bucket,
		client: client,
	}
}
