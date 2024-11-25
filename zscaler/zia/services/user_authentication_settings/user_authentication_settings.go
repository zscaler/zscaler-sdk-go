package user_authentication_settings

import (
	"context"
	"fmt"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	exemptedUrlsEndpoint = "/zia/api/v1/authSettings/exemptedUrls"
)

type ExemptedUrls struct {
	// Domains or URLs which are exempted from SSL Inspection
	URLs []string `json:"urls"`
}

type QueryParameters struct {
	ID string
}

func Get(ctx context.Context, service *zscaler.Service) (*ExemptedUrls, error) {
	var urls ExemptedUrls
	err := service.Client.Read(ctx, exemptedUrlsEndpoint, &urls)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning exempted url rules from Get: %v", urls)
	return &urls, nil
}

// return the new items that were added to slice1.
func difference(slice1 []string, slice2 []string) []string {
	var diff []string
	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, s1)
		}
	}
	return diff
}

func Update(ctx context.Context, service *zscaler.Service, urls ExemptedUrls) (*ExemptedUrls, error) {
	currentUrsl, err := Get(ctx, service)
	if err != nil {
		return nil, err
	}
	newUrls := difference(urls.URLs, currentUrsl.URLs)
	removedUrls := difference(currentUrsl.URLs, urls.URLs)
	if len(newUrls) > 0 {
		_, err := service.Client.Create(ctx, fmt.Sprintf("%s?action=ADD_TO_LIST", exemptedUrlsEndpoint), ExemptedUrls{newUrls})
		if err != nil {
			return nil, err
		}
	}
	if len(removedUrls) > 0 {
		_, err := service.Client.Create(ctx, fmt.Sprintf("%s?action=REMOVE_FROM_LIST", exemptedUrlsEndpoint), ExemptedUrls{removedUrls})
		if err != nil {
			return nil, err
		}
	}
	return &urls, nil
}
