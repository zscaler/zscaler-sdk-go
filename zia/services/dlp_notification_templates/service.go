package dlp_notification_templates

import (
	"github.com/zscaler/zscaler-sdk-go/zia"
)

type Service struct {
	Client *zia.Client
}

func New(c *zia.Client) *Service {
	return &Service{Client: c}
}
