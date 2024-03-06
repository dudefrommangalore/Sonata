package awslib

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"google.golang.org/appengine/log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	signer    *v4.Signer
	tokenOnce sync.Once
	profile   = flag.String("profile", "presto/dev", "Profile to use")
	region    = flag.String("region", "us-west-2", "AWS region to use")
)

func NewSigner(ctx context.Context) (*v4.Signer, error) {
	tokenOnce.Do(func() {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
		sess, err := session.NewSessionWithOptions(session.Options{
			Profile:            *profile,
			Config:             aws.Config{Region: aws.String(*region)},
			AssumeRoleDuration: time.Duration(24) * time.Hour,
		})
		if err == nil {
			signer = v4.NewSigner(sess.Config.Credentials)
		} else {
			log.Errorf(ctx, "failed to create new session %v", err)
		}
	})
	if signer == nil {
		return nil, fmt.Errorf("Signer is nil")
	}

	return signer, nil
}

func SignRequest(ctx context.Context, req *http.Request, service string) error {
	signer, err := NewSigner(ctx)
	if err != nil {
		return err
	}

	signer.Sign(req, nil, service, *region, time.Now())
	return nil
}
