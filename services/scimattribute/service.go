package scimattribute

import (
	"github.com/willguibr/zscaler-sdk-go/zpa"
	"github.com/willguibr/zscaler-sdk-go/zpa/client"
)

type Service struct {
	Client *client.Client
}

func New(c *zpa.Config) *Service {
	return &Service{Client: client.NewClient(c)}
}
