package urlcategories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	urlCategoriesEndpoint = "/zia/api/v1/urlCategories"
	urlQuotaHandler       = "/urlQuota"
	urlLookupEndpoint     = "/zia/api/v1/urlLookup"
)

type URLCategory struct {
	// URL category
	ID string `json:"id,omitempty"`

	// Name of the URL category. This is only required for custom URL categories.
	ConfiguredName string `json:"configuredName,omitempty"`

	// Custom keywords associated to a URL category. Up to 2048 custom keywords can be added per organization across all categories (including bandwidth classes).
	Keywords []string `json:"keywords"`

	// Retained custom keywords from the parent URL category that is associated to a URL category. Up to 2048 retained parent keywords can be added per organization across all categories (including bandwidth classes).
	KeywordsRetainingParentCategory []string `json:"keywordsRetainingParentCategory"`

	// Custom URLs to add to a URL category. Up to 25,000 custom URLs can be added per organization across all categories (including bandwidth classes).
	Urls []string `json:"urls"`

	// URLs added to a custom URL category are also retained under the original parent URL category (i.e., the predefined category the URL previously belonged to). The URLs entered are covered by policies that reference the original parent URL category as well as those that reference the custom URL category. For example, if you add www.amazon.com, this URL is covered by policies that reference the custom URL category as well as policies that reference its parent URL category of "Online Shopping".
	DBCategorizedUrls []string `json:"dbCategorizedUrls"`

	//
	CustomCategory bool `json:"customCategory"`

	// Scope of the custom categories.
	Scopes []Scopes `json:"scopes,omitempty"`

	// Value is set to false for custom URL category when due to scope user does not have edit permission
	Editable bool `json:"editable"`

	// Description of the URL category. Contains tag name and needs to be localized on client side in case of predefined category (customCategory=null or =false), else it contains the user-entered description which does not have localization support.
	Description string `json:"description,omitempty"`

	// Type of the custom categories.
	Type string `json:"type,omitempty"`

	// URL and keyword counts for the URL category.
	URLKeywordCounts *URLKeywordCounts `json:"urlKeywordCounts,omitempty"`
	Val              int               `json:"val,omitempty"`

	// The number of custom URLs associated to the URL category.
	CustomUrlsCount int `json:"customUrlsCount,omitempty"`

	// Super Category of the URL category. This field is required when creating custom URL categories.
	SuperCategory string `json:"superCategory,omitempty"`

	// Specifies the type of URL match, such as using exact URLs or regex patterns.
	// For regex, the patterns can be specified using the regexPatterns and regexPatternsRetainingParentCategory fields.
	// For exact URLs, specify the required URLs, keywords, or IP ranges using the appropriate fields.
	// Supported Values: EXACT, REGEX
	// Note: To enable the Regex feature, contact Zscaler Support.
	UrlType string `json:"urlType,omitempty"`

	// The number of custom URLs associated to the URL category, that also need to be retained under the original parent category.
	UrlsRetainingParentCategoryCount int `json:"urlsRetainingParentCategoryCount"`

	// Custom IP address ranges associated to a URL category. Up to 2000 custom IP address ranges and retaining parent custom IP address ranges can be added, per organization, across all categories.
	IPRanges []string `json:"ipRanges"`

	// The retaining parent custom IP address ranges associated to a URL category. Up to 2000 custom IP ranges and retaining parent custom IP address ranges can be added, per organization, across all categories.
	IPRangesRetainingParentCategory []string `json:"ipRangesRetainingParentCategory"`

	// The number of custom IP address ranges associated to the URL category.
	CustomIpRangesCount int `json:"customIpRangesCount"`

	// The number of custom IP address ranges associated to the URL category, that also need to be retained under the original parent category.
	IPRangesRetainingParentCategoryCount int `json:"ipRangesRetainingParentCategoryCount"`

	// Regex patterns associated with this custom category to support full or partial matches with URLs. Multiple patterns can be added to a category. A pattern must be between 3 and 255 characters, and a maximum of 256 patterns can be added across all categories.
	RegexPatterns []string `json:"regexPatterns"`

	// This field specifies regex patterns for URLs that must be covered by policies directly referencing this custom category as well as by policies referencing its parent URL category (e.g., Corporate Marketing). Multiple patterns can be added to a category.
	// A pattern must be between 3 and 255 characters, and a maximum of 256 patterns can be added across all categories.
	RegexPatternsRetainingParentCategory []string `json:"regexPatternsRetainingParentCategory"`
}

type Scopes struct {
	// Only applicable for the LOCATION_GROUP admin scope type, in which case this attribute gives the list of ID/name pairs of locations within the location group. The attribute name is subject to change.
	ScopeGroupMemberEntities []common.IDNameExtensions `json:"scopeGroupMemberEntities,omitempty"`

	// The admin scope type. The attribute name is subject to change.
	Type string `json:"Type,omitempty"`

	// Based on the admin scope type, the entities can be the ID/name pair of departments, locations, or location groups. The attribute name is subject to change.
	ScopeEntities []common.IDNameExtensions `json:"ScopeEntities,omitempty"`
}

type URLKeywordCounts struct {
	// Custom URL count for the category.
	TotalURLCount int `json:"totalUrlCount,omitempty"`

	// Count of URLs with retain parent category.
	RetainParentURLCount int `json:"retainParentUrlCount,omitempty"`

	// Total keyword count for the category.
	TotalKeywordCount int `json:"totalKeywordCount,omitempty"`

	// Count of total keywords with retain parent category.
	RetainParentKeywordCount int `json:"retainParentKeywordCount,omitempty"`
}

type URLQuota struct {
	UniqueUrlsProvisioned int `json:"uniqueUrlsProvisioned,omitempty"`
	RemainingUrlsQuota    int `json:"remainingUrlsQuota,omitempty"`
}

type URLClassification struct {
	URL                                 string   `json:"url,omitempty"`
	URLClassifications                  []string `json:"urlClassifications,omitempty"`
	URLClassificationsWithSecurityAlert []string `json:"urlClassificationsWithSecurityAlert,omitempty"`
	Application                         string   `json:"application,omitempty"`
}

type DomainMatch struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type URLReview struct {
	URL        string        `json:"url"`
	DomainType string        `json:"domainType"`
	Matches    []DomainMatch `json:"matches"`
}

func Get(ctx context.Context, service *zscaler.Service, categoryID string) (*URLCategory, error) {
	var urlCategory URLCategory
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%s", urlCategoriesEndpoint, categoryID), &urlCategory)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning custom url category from Get: %s", urlCategory.ID)
	return &urlCategory, nil
}

func GetCustomURLCategories(ctx context.Context, service *zscaler.Service, customName string, includeOnlyUrlKeywordCounts, customOnly bool, categoryType string) (*URLCategory, error) {
	var urlCategory []URLCategory
	queryParams := url.Values{}

	if includeOnlyUrlKeywordCounts {
		queryParams.Set("includeOnlyUrlKeywordCounts", "false")
	}
	if customOnly {
		queryParams.Set("customOnly", "true")
	}
	// Add type parameter to filter by category type (ALL, URL_CATEGORY, TLD_CATEGORY)
	if categoryType != "" {
		queryParams.Set("type", categoryType)
	}

	err := service.Client.Read(ctx, fmt.Sprintf("%s?%s", urlCategoriesEndpoint, queryParams.Encode()), &urlCategory)
	if err != nil {
		return nil, err
	}

	for _, custom := range urlCategory {
		if strings.EqualFold(custom.ConfiguredName, customName) {
			return &custom, nil
		}
	}
	return nil, fmt.Errorf("no custom url category found with name: %s", customName)
}

func GetAllCustomURLCategories(ctx context.Context, service *zscaler.Service) ([]URLCategory, error) {
	var all []URLCategory
	queryParams := url.Values{}
	queryParams.Set("customOnly", "true")

	err := service.Client.Read(ctx, fmt.Sprintf("%s?%s", urlCategoriesEndpoint, queryParams.Encode()), &all)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func CreateURLCategories(ctx context.Context, service *zscaler.Service, category *URLCategory) (*URLCategory, error) {
	resp, err := service.Client.Create(ctx, urlCategoriesEndpoint, *category)
	if err != nil {
		return nil, err
	}

	createdUrlCategory, ok := resp.(*URLCategory)
	if !ok {
		return nil, errors.New("object returned from API was not a url category Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning url category from Create: %v", createdUrlCategory.ID)
	return createdUrlCategory, nil
}

func UpdateURLCategories(ctx context.Context, service *zscaler.Service, categoryID string, category *URLCategory, action string) (*URLCategory, *http.Response, error) {
	endpoint := fmt.Sprintf("%s/%s", urlCategoriesEndpoint, categoryID)

	// Append action query parameter if provided
	if action != "" {
		endpoint = fmt.Sprintf("%s?action=%s", endpoint, url.QueryEscape(action))
	}

	resp, err := service.Client.UpdateWithPut(ctx, endpoint, *category)
	if err != nil {
		return nil, nil, err
	}
	updatedUrlCategory, _ := resp.(*URLCategory)
	service.Client.GetLogger().Printf("[DEBUG]Returning url category from Update: %s", updatedUrlCategory.ID)
	return updatedUrlCategory, nil, nil
}

func DeleteURLCategories(ctx context.Context, service *zscaler.Service, categoryID string) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%s", urlCategoriesEndpoint, categoryID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetURLQuota(ctx context.Context, service *zscaler.Service) (*URLQuota, error) {
	url := fmt.Sprintf("%s/%s", urlCategoriesEndpoint, urlQuotaHandler)
	var quota URLQuota
	err := service.Client.Read(ctx, url, &quota)
	if err != nil {
		return nil, err
	}
	return &quota, nil
}

func GetURLLookup(ctx context.Context, service *zscaler.Service, urls []string) ([]URLClassification, error) {
	resp, err := service.Client.CreateWithSlicePayload(ctx, urlLookupEndpoint, urls)
	if err != nil {
		return nil, err
	}

	var lookupResults []URLClassification
	err = json.Unmarshal(resp, &lookupResults)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] returning URL lookup results: %+v", lookupResults)
	return lookupResults, nil
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]URLCategory, error) {
	var urlCategories []URLCategory
	err := common.ReadAllPages(ctx, service.Client, urlCategoriesEndpoint+"/lite", &urlCategories)
	if err != nil {
		service.Client.GetLogger().Printf("[ERROR] Error fetching URL categories: %v", err)
		return nil, err
	}
	return urlCategories, nil
}

func CreateURLReview(ctx context.Context, service *zscaler.Service, domains []string) ([]URLReview, error) {
	resp, err := service.Client.CreateWithSlicePayload(ctx, urlCategoriesEndpoint+"/review/domains", domains)
	if err != nil {
		return nil, err
	}

	var reviewResults []URLReview
	err = json.Unmarshal(resp, &reviewResults)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] returning URL review results: %+v", reviewResults)
	return reviewResults, nil
}

func UpdateURLReview(ctx context.Context, service *zscaler.Service, reviews []URLReview) error {
	resp, err := service.Client.UpdateWithSlicePayload(context.Background(), urlCategoriesEndpoint+"/review/domains", reviews)
	if err != nil {
		return err
	}

	if len(resp) > 0 {
		return errors.New("unexpected response format")
	}

	service.Client.GetLogger().Printf("[DEBUG] successfully updated URL review")
	return nil
}

func GetAll(ctx context.Context, service *zscaler.Service, customOnly, includeOnlyUrlKeywordCounts bool, categoryType string) ([]URLCategory, error) {
	var urlCategories []URLCategory

	// Build the endpoint with optional query parameters
	endpoint := urlCategoriesEndpoint
	queryParams := url.Values{}

	if customOnly {
		queryParams.Set("customOnly", "true")
	}
	if includeOnlyUrlKeywordCounts {
		queryParams.Set("includeOnlyUrlKeywordCounts", "true")
	}
	// Add type parameter to filter by category type (ALL, URL_CATEGORY, TLD_CATEGORY)
	// If categoryType is empty, no type filter is applied (returns predefined + custom URL_CATEGORY)
	// If categoryType is "ALL", returns predefined + custom categories of all types
	if categoryType != "" {
		queryParams.Set("type", categoryType)
	}

	// Append query parameters to endpoint if any exist
	if len(queryParams) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())
	}

	// Use service.Client.Read directly since the API doesn't support pagination
	// The API returns all results in a single response
	err := service.Client.Read(ctx, endpoint, &urlCategories)
	return urlCategories, err
}
