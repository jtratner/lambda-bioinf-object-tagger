package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"log"
	"regexp"
)

const FILETYPE_KEY = "filetype"

var Version string
var GitCommit string

var REGEXES = map[string]*regexp.Regexp{
	"fastq": regexp.MustCompile(".*.fastq(.gz)?$"),
	"bam":   regexp.MustCompile(".*.bam$"),
}

func entityPath(obj *events.S3Entity) string {
	return fmt.Sprintf("s3://%s/%s", obj.Bucket.Name, obj.Object.Key)
}

func init() {
	log.SetPrefix(fmt.Sprintf("Version:%s|Commit:%s", Version, GitCommit))
}

type LambdaResponse struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}

func LambdaHandler(ctx context.Context, evt *events.S3Event) (*LambdaResponse, error) {
	tagsApplied := 0
	for _, rec := range evt.Records {
		tagForObject := getTagForObject(&rec.S3)
		if tagForObject != nil {
			if sess, err := session.NewSession(); err != nil {
				return nil, err
			} else {
				output, err := applyTagToObject(tagForObject, s3.New(sess))
				if err != nil {
					return nil, err
				}
				log.Printf("successfully applied tag to %s (%s)", entityPath(&rec.S3), output.String())
				tagsApplied++
			}
		}
	}
	return &LambdaResponse{Count: tagsApplied, Message: "completed successfully"}, nil
}

func applyTagToObject(inpt *s3.PutObjectTaggingInput, svc s3iface.S3API) (*s3.PutObjectTaggingOutput, error) {
	return svc.PutObjectTagging(inpt)
}

func getTagForObject(obj *events.S3Entity) *s3.PutObjectTaggingInput {
	filetypeKey := FILETYPE_KEY
	for filetype, regex := range REGEXES {
		if regex.MatchString(obj.Object.Key) {
			return &s3.PutObjectTaggingInput{
				Bucket: &obj.Bucket.Name,
				Key:    &obj.Object.Key,
				Tagging: &s3.Tagging{
					TagSet: []*s3.Tag{
						{Key: &filetypeKey, Value: &filetype},
					},
				},
				VersionId: &obj.Object.VersionID,
			}
		}
	}
	return nil

}

func main() {
	lambda.Start(LambdaHandler)
}
