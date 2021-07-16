package aws

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/joho/godotenv"
)

func AmazonSES(subject, body string) {

	godotenv.Load()
	AccessKeyId := os.Getenv("ses_accessKeyId")
	SecretAccessKey := os.Getenv("ses_secretAccessKey")
	MyRegion := os.Getenv("region")

	awsSession := session.New(&aws.Config{
        Region:      aws.String(MyRegion),
        Credentials: credentials.NewStaticCredentials(AccessKeyId, SecretAccessKey , ""),
    })

sesSession := ses.New(awsSession)

sesEmailInput := &ses.SendEmailInput{
    Destination: &ses.Destination{
        ToAddresses:  []*string{aws.String("send_email")},
    },
    Message: &ses.Message{
        Body: &ses.Body{
            Text: &ses.Content{
                Data: aws.String(body)},
        },
        Subject: &ses.Content{
            Data: aws.String(subject),
        },
    },
    Source: aws.String("from_email"),

}

	sesSession.SendEmail(sesEmailInput)
}