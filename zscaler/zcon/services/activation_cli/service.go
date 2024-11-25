package main

import "github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zcon"

type Service struct {
	Client *zcon.Client
}

func New(c *zcon.Client) *Service {
	return &Service{Client: c}
}
