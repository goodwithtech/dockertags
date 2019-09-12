package ecr

import (
	"context"
	"strings"
	"time"

	"github.com/goodwithtech/dockertags/pkg/log"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/goodwithtech/dockertags/pkg/types"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	service "github.com/aws/aws-sdk-go/service/ecr"
)

type ECR struct{}

var _ time.Duration
var _ strings.Reader
var _ aws.Config

func getSession(option types.RequestOption) (*session.Session, error) {
	// create custom credential information if option is valid
	if option.AwsSecretKey != "" && option.AwsAccessKey != "" && option.AwsRegion != "" {
		return session.NewSessionWithOptions(
			session.Options{
				Config: aws.Config{
					Region: aws.String(option.AwsRegion),
					Credentials: credentials.NewStaticCredentialsFromCreds(
						credentials.Value{
							AccessKeyID:     option.AwsAccessKey,
							SecretAccessKey: option.AwsSecretKey,
						},
					),
				},
			})
	}
	// use shared configuration normally
	return session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
}

func (p *ECR) Run(ctx context.Context, domain, repository string, option types.RequestOption) (types.ImageTags, error) {
	sess, err := getSession(option)
	if err != nil {
		return nil, err
	}
	svc := service.New(sess)
	input := &service.DescribeImagesInput{
		RepositoryName: aws.String(repository),
		// Only show tagged image
		Filter: &service.DescribeImagesFilter{TagStatus: aws.String("TAGGED")},
	}

	result, err := svc.DescribeImages(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case service.ErrCodeServerException:
				log.Logger.Errorf(service.ErrCodeServerException, aerr.Error())
			case service.ErrCodeInvalidParameterException:
				log.Logger.Errorf(service.ErrCodeInvalidParameterException, aerr.Error())
			case service.ErrCodeRepositoryNotFoundException:
				log.Logger.Errorf(service.ErrCodeRepositoryNotFoundException, aerr.Error())
			default:
				log.Logger.Errorf(aerr.Error())
			}
		} else {
			log.Logger.Error(err.Error())
		}
		return nil, err
	}

	imageTags := []types.ImageTag{}
	for _, detail := range result.ImageDetails {
		if len(detail.ImageTags) == 0 {
			continue
		}
		tags := []string{}
		for _, t := range detail.ImageTags {
			if t != nil {
				tags = append(tags, *t)
			}
		}

		imageTags = append(imageTags, types.ImageTag{
			Tags:       tags,
			Byte:       getIntByte(detail.ImageSizeInBytes),
			CreatedAt:  nil,
			UploadedAt: detail.ImagePushedAt,
		})
	}

	return imageTags, nil
}

func getIntByte(b *int64) *int {
	if b == nil {
		return nil
	}
	i := int(*b)
	return &i
}
