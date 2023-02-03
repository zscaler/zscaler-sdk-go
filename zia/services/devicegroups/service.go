package devicegroups

import (
	"github.com/zscaler/zscaler-sdk-go/v1/zia"
)

type Service struct {
	Client *zia.Client
}

func New(c *zia.Client) *Service {
	return &Service{Client: c}
}
