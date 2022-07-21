package provisioningkey

import (
	"github.com/willguibr/zscaler-sdk-go/zpa"
)

type Service struct {
	Client *zpa.Client
}

func New(c *zpa.Client) *Service {
	return &Service{Client: c}
}
