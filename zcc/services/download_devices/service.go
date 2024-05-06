package download_devices

import (
	"github.com/zscaler/zscaler-sdk-go/v2/zcc"
)

type Service struct {
	Client *zcc.Client
}

func New(c *zcc.Client) *Service {
	return &Service{Client: c}
}
