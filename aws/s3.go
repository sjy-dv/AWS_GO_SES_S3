package aws

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

func ConnS3() *session.Session {
	godotenv.Load()
	AccessKeyId := os.Getenv("accessKeyId")
	SecretAccessKey := os.Getenv("secretAccessKey")
	MyRegion := os.Getenv("region")

	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(AccessKeyId, SecretAccessKey, ""),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	return sess
}

var sess = ConnS3()

func SingleUploadFile(r *http.Request) string {

	// form-data에 img라는 파일이 있다고 가정
	// Assume there is a file called img in form-data
	file, header, _ := r.FormFile("img")

	var extension string
	var err error
	defer file.Close()

	filename := header.Filename // if you want to more security + alpha to filename

	// not using this , you r not using image file in mobile app application
	if filepath.Ext(filename) == ".jpg" {
		extension = "image/jpg"
	} else if filepath.Ext(filename) == ".jpeg" {
		extension = "image/jpeg"
	} else if filepath.Ext(filename) == ".png" {
		extension = "image/png"
	} else if filepath.Ext(filename) == ".gif" {
		extension = "image/gif"
	}

	url := "your s3 endpointurl"

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("your bucket"),
		Key: aws.String(filename),
		ACL: aws.String("public-read"),
		Body: file,
		ContentType: aws.String(extension),
	})
	if err != nil {
		log.Fatal(err)
	}
	save_in_database_filename := url + "/bucketaddress/" + filename
	return save_in_database_filename	
}

func SingleDeleteFile(deleteimagename string) error {

	var err error

	deleter := s3.New(sess)
	filter_in_database_filename := strings.Replace(deleteimagename, "your bucket address","",1)

	_, err = deleter.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: aws.String("/bucket file address"),
			Key : aws.String(filter_in_database_filename),
		},
	)
	
	if err != nil {
		return err
	}
	return nil
}

func MultiFileUpload(r *http.Request) (multi_image_set string){

	var err error

	m := r.MultipartForm

	var image_set []string
	files := m.File["imgs"]

	var extension string

	for i := range files {

		file, _ := files[i].Open()
		defer file.Close()

		filename := files[i].Filename // + alpha some thing in your working project

		if filepath.Ext(filename) == ".jpg" {
			extension = "image/jpg"
		} else if filepath.Ext(filename) == ".jpeg" {
			extension = "image/jpeg"
		} else if filepath.Ext(filename) == ".png" {
			extension = "image/png"
		} else if filepath.Ext(filename) == ".gif" {
			extension = "image/gif"
		}

		url := "your s3 address"
		uploader := s3manager.NewUploader(sess)

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String("your bucket"),
			Key:    aws.String(filename),
			ACL:    aws.String("public-read"),
			Body:   file,
			ContentType:  aws.String(extension),
		})
		
		if err != nil {
			log.Fatal(err)
		}

		each_filename := url + "/your bucket/" + filename
		image_set = append(image_set, each_filename)
	}

	multi_image := strings.Join(image_set, ",")
	
	return multi_image
}

func MultiFileDeleter(delete_image_set string) error {

	var err error

	var files []string
	files = strings.Split(delete_image_set, ",")
	deleter := s3.New(sess)

	for i := range files {
	filter_in_database_filename := strings.Replace(files[i],"your s3 bucket address","",1)
	
	_, err = deleter.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: aws.String("your bucket"), 
			Key: aws.String(filter_in_database_filename),
		},)	
	return err
	}
	return nil
}
