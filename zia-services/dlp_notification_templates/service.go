package dlp_notification_templates

import (
	"github.com/willguibr/zscaler-sdk-go/zia/client"
)

type Service struct {
	Client *client.Client
}

func New(c *client.Client) *Service {
	return &Service{Client: c}
}
