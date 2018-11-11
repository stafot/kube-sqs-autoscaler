package sqs

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"strconv"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/pkg/errors"
)

type SQS interface {
	GetQueueAttributes(*sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error)
	// only implemented on unit tests
	SetQueueAttributes(*sqs.SetQueueAttributesInput) (*sqs.SetQueueAttributesOutput, error)
}

type SqsClient struct {
	Client   SQS
	QueueUrl string
}

func NewSqsClient(queue string, region string) *SqsClient {
	key := os.Getenv("SQS_ACCESS_KEY")
	secret := os.Getenv("SQS_ACCESS_SECRET")
	if (key != "" && secret != "") {
		svc := sqs.New(session.New(), &aws.Config{Region: aws.String(region), Credentials: credentials.NewStaticCredentials(key, secret, "")})
		return &SqsClient{
			svc,
			queue,
		}
	} else {
		svc := sqs.New(session.New(), &aws.Config{Region: aws.String(region)})
		return &SqsClient{
			svc,
			queue,
		}
	}
}

func (s *SqsClient) NumMessages() (int, error) {
	params := &sqs.GetQueueAttributesInput{
		AttributeNames: []*string{aws.String("ApproximateNumberOfMessages")},
		QueueUrl:       aws.String(s.QueueUrl),
	}

	out, err := s.Client.GetQueueAttributes(params)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to get messages in SQS")
	}

	messages, err := strconv.Atoi(*out.Attributes["ApproximateNumberOfMessages"])
	if err != nil {
		return 0, errors.Wrap(err, "Failed to get number of messages in queue")
	}

	return messages, nil
}
