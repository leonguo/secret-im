package controllers
//
//import (
//	"github.com/labstack/echo"
//	"net/http"
//	"../../util"
//	"github.com/aws/aws-sdk-go/aws/session"
//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/service/s3"
//	"time"
//	"log"
//	"../../config"
//	"github.com/aws/aws-sdk-go/aws/credentials"
//	"strconv"
//)
//
//func GetAttachment(c echo.Context) error {
//	attachmentId, _ := util.GenerateAttachmentId()
//	attachmentIdString := strconv.FormatInt(attachmentId, 10)
//	creds := credentials.NewStaticCredentials(config.AppConfig.GetString("s3_attachments.accessKey"), config.AppConfig.GetString("s3_attachments.accessSecret"), "")
//	sess, err := session.NewSession(&aws.Config{
//		Region:      aws.String(config.AppConfig.GetString("s3_attachments.region")),
//		Credentials: creds},
//	)
//	// Create S3 service client
//	svc := s3.New(sess)
//	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
//		Bucket: aws.String(config.AppConfig.GetString("s3_attachments.bucket")),
//		Key:    aws.String(attachmentIdString),
//	})
//	urlStr, err := req.Presign(15 * time.Minute)
//	if err != nil {
//		log.Println("Failed to sign request", err)
//		return c.JSON(http.StatusInternalServerError, err)
//	}
//	result := echo.Map{"id": attachmentId, "idString": attachmentIdString, "location": urlStr}
//	return c.JSON(http.StatusOK, result)
//}
//
//func GetAttachmentId(c echo.Context) error {
//	attachmentIdString := c.Param("attachmentId")
//	creds := credentials.NewStaticCredentials(config.AppConfig.GetString("s3_attachments.accessKey"), config.AppConfig.GetString("s3_attachments.accessSecret"), "")
//	sess, err := session.NewSession(&aws.Config{
//		Region:      aws.String(config.AppConfig.GetString("s3_attachments.region")),
//		Credentials: creds},
//	)
//	svc := s3.New(sess)
//	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
//		Bucket: aws.String(config.AppConfig.GetString("s3_attachments.bucket")),
//		Key:    aws.String(attachmentIdString),
//	})
//	urlStr, err := req.Presign(15 * time.Minute)
//	if err != nil {
//		log.Println("Failed to sign request", err)
//		return c.JSON(http.StatusInternalServerError, err)
//	}
//	result := echo.Map{"location": urlStr}
//	return c.JSON(http.StatusOK, result)
//}
