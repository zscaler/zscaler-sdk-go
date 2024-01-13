package dlp_exact_data_match_lite

import (
	"github.com/zscaler/zscaler-sdk-go/v2/zia"
)

type Service struct {
	Client *zia.Client
}

func New(c *zia.Client) *Service {
	return &Service{Client: c}
}
