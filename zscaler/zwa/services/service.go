package services

import (
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa"
)

type Service struct {
	Client *zwa.Client
}

func New(c *zwa.Client) *Service {
	return &Service{Client: c}
}
