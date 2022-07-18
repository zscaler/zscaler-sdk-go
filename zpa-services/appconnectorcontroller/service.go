package appconnectorcontroller

import (
	"github.com/willguibr/zscaler-sdk-go/client"
	"github.com/willguibr/zscaler-sdk-go/gozscaler"
)

type Service struct {
	Client *client.Client
}

func New(c *gozscaler.Config) *Service {
	return &Service{Client: client.NewClient(c)}
}
