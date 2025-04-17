package services

import "github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw"

type Service struct {
	Client *ztw.Client
}

func New(c *ztw.Client) *Service {
	return &Service{Client: c}
}
