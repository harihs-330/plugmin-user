package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSProvider struct {
	UnImplemented
	Cfg    *aws.Config
	bucket Storage
}

var _ Provider = (*AWSProvider)(nil)

func (srv *AWSProvider) Load(_ *Config) Provider {
	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil
	}
	srv.Cfg = &awsCfg
	srv.bucket = &Bucket{client: s3.NewFromConfig(awsCfg)}

	return srv
}

func (srv *AWSProvider) Storage() Storage {
	return srv.bucket
}

type Bucket struct {
	client *s3.Client
}

func (bucket *Bucket) Get(string) ([]byte, error) {
	return []byte{}, nil
}

func (bucket *Bucket) Put(string, []byte) (bool, error) {
	return false, nil
}
