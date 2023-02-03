package customerversionprofile

import (
	"github.com/zscaler/zscaler-sdk-go/v1/zpa"
)

type Service struct {
	Client *zpa.Client
}

func New(c *zpa.Client) *Service {
	return &Service{Client: c}
}
