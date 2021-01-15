package helpers

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"net/http"
	"sort"
	"time"
)

func SaveInS3(s *session.Session, s3Bucket string, fileName string, body []byte) error {
	r, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s3Bucket),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(body),
		ContentLength:        aws.Int64(int64(len(body))),
		ContentType:          aws.String(http.DetectContentType(body)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		CacheControl:         aws.String("no-cache"),
	})
	if err != nil {
		log.Printf("Error saving in s3: %s", fileName)
		fmt.Println(err)
	} else {
		log.Printf("Uploaded file: %s", r.String())
	}
	return err
}

func DeleteFromS3(sess *session.Session, bucket string, key string) bool {
	deleteObjectOutput, err := s3.New(sess).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return false
	}
	return deleteObjectOutput.VersionId != nil
}

func DownloadFromS3(sess *session.Session, bucket string, key string) []byte {
	downloader := s3manager.NewDownloader(sess)
	latestUpdateBuff := &aws.WriteAtBuffer{}
	downloaded, err := downloader.Download(latestUpdateBuff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		log.Printf("Error while downloading bytes (%d) from %s -> %s, %v", downloaded, bucket, key, err.Error())
		return []byte{}
	}
	dled := latestUpdateBuff.Bytes()
	log.Printf("Downloaded %d bytes from %s -> %s::%s", downloaded, bucket, key, dled)
	return dled
}

func ListAllObjectFromPrefix(sess *session.Session, bucket string, prefix string) []*s3.Object {
	var objectList []*s3.Object
	err := s3.New(sess).ListObjectsPages(&s3.ListObjectsInput{
		Bucket: &bucket,
		Prefix: aws.String(prefix),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, object := range p.Contents {
			objectList = append(objectList, object)
		}
		return true
	})
	if err != nil {
		log.Panic(err)
	}
	return objectList
}

func DownloadAllObjectsFromPrefix(sess *session.Session, bucket string, prefix string) map[string]string {
	objects := ListAllObjectFromPrefix(sess, bucket, prefix)
	downloaded := make(map[string]string)
	for _, object := range objects {
		fromS3 := DownloadFromS3(sess, bucket, *object.Key)
		downloaded[*object.Key] = string(fromS3)
	}
	return downloaded
}

func ObjectExists(sess *session.Session, bucket string, key string) bool {
	svc := s3.New(sess)
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return false
	}
	return true
}

func GetLatestObjectFromPrefix(sess *session.Session, bucket string, prefix string) *s3.Object {
	var objectList []*s3.Object
	err := s3.New(sess).ListObjectsPages(&s3.ListObjectsInput{
		Bucket: &bucket,
		Prefix: &prefix,
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, object := range p.Contents {
			objectList = append(objectList, object)
		}
		return false
	})
	if err != nil {
		log.Panic(err)
	}
	sort.SliceStable(objectList, func(i, j int) bool {
		return objectList[i].LastModified.Unix() < objectList[j].LastModified.Unix()
	})
	log.Printf("Found objects count: %d\n", len(objectList))
	if len(objectList) == 0 {
		return nil
	}
	return objectList[len(objectList)-1]
}

func PresignedDownloadUrl(sess *session.Session, bucket string, key string, duration time.Duration) string {
	req, _ := s3.New(sess).GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String("attachment;filename=client-inventory"),
		ResponseContentType:        aws.String("application/octet-stream"),
		ResponseExpires:            aws.Time(time.Now().Add(duration)),
	})
	urlStr, err := req.Presign(duration)
	if err != nil {
		log.Fatal(err)
	}
	return urlStr
}

func PresignedUploadUrl(sess *session.Session, bucket string, key string, filename string, duration time.Duration) string {
	request, _ := s3.New(sess).PutObjectRequest(&s3.PutObjectInput{
		Bucket:             aws.String(bucket),
		CacheControl:       aws.String("no-cache"),
		ContentDisposition: aws.String(filename),
		Key:                aws.String(key),
	})
	presign, err := request.Presign(duration)
	if err != nil {
		log.Fatal(err)
	}
	return presign
}
