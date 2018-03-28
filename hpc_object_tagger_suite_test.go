package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getEntity(bucket string, key string, size int) *events.S3Entity {
	return &events.S3Entity{
		Bucket: events.S3Bucket{Name: bucket},
		Object: events.S3Object{Key: key, Size: int64(size)},
	}
}

func getEvent(bucket string, key string, size int) *events.S3Event {
	return &events.S3Event{
		Records: []events.S3EventRecord{events.S3EventRecord{S3: *getEntity(bucket, key, size)}},
	}
}

func TestLambdaS3Tagger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HPCObjectTagger Suite")
}

type mockPutObjectTagging struct {
	s3iface.S3API
	Resp  s3.PutObjectTaggingOutput
	Error error
}

func (m *mockPutObjectTagging) PutObjectTagging(*s3.PutObjectTaggingInput) (*s3.PutObjectTaggingOutput, error) {
	return &m.Resp, m.Error
}

func (m *mockPutObjectTagging) PutObjectTaggingWithContext(aws.Context, *s3.PutObjectTaggingInput, ...request.Option) (*s3.PutObjectTaggingOutput, error) {
	return &m.Resp, m.Error
}

func expectTag(result *s3.PutObjectTaggingInput, key string, value string, gingkoParams ...interface{}) {
	Expect(result).ToNot(BeNil(), gingkoParams...)
	Expect(len(result.Tagging.TagSet)).To(Equal(1), gingkoParams...)
	Expect(*result.Tagging.TagSet[0].Key).To(Equal(key), gingkoParams...)
	Expect(*result.Tagging.TagSet[0].Value).To(Equal(value), gingkoParams...)
}

var _ = Describe("lambda tagger", func() {
	Context("entity to Path", func() {
		It("should work", func() {
			Expect(entityPath(getEntity("mybucket", "mykey.gz", 1))).To(Equal("s3://mybucket/mykey.gz"))
		})
	})
	Context("getTagForObject", func() {
		It("should only return tag if matches", func() {
			Expect(getTagForObject(getEntity("bucket", "key", 1))).To(BeNil())
		})
		It("should return tag on match", func() {
			result := getTagForObject(getEntity("bucket", "mykey.fastq.gz", 50*MB))
			expectTag(result, "filetype", "fastq")
			Expect(*result.Bucket).To(Equal("bucket"))
		})
		It("should tag really big files as largefile", func() {
			result := getTagForObject(getEntity("bucket", "mykey.other", 50*MB))
			expectTag(result, "filetype", "largefile")
			Expect(*result.Bucket).To(Equal("bucket"))
		})
	})
	Context("handleEvent", func() {
		It("should work with a mocked object", func() {
			m := &mockPutObjectTagging{}
			m.PutObjectTagging(nil)
			Expect(handleEvent(context.Background(), getEvent("bucket", "my.fastq.gz", 1), m)).To(Equal(
				&LambdaResponse{1, "completed successfully"},
			))
		})
	})
})
