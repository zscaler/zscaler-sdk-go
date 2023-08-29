package urlcategories

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
)

const (
	urlCategoriesEndpoint = "/urlCategories"
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

func (service *Service) Get(categoryID string) (*URLCategory, error) {
	var urlCategory URLCategory
	err := service.Client.Read(fmt.Sprintf("%s/%s", urlCategoriesEndpoint, categoryID), &urlCategory)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning custom url category from Get: %s", urlCategory.ID)
	return &urlCategory, nil
}

func (service *Service) GetCustomURLCategories(customName string) (*URLCategory, error) {
	var urlCategory []URLCategory
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?customOnly=%s", urlCategoriesEndpoint, "true"), &urlCategory)
	if err != nil {
		return nil, err
	}
	for _, custom := range urlCategory {
		if strings.EqualFold(custom.ConfiguredName, customName) { // Use ConfiguredName instead of ID for comparison
			return &custom, nil
		}
	}
	return nil, fmt.Errorf("no custom url category found with name: %s", customName)
}

func (service *Service) GetIncludeOnlyUrlKeyWordCounts(customName string) (*URLCategory, error) {
	var urlCategory []URLCategory
	err := service.Client.Read(fmt.Sprintf("%s?includeOnlyUrlKeywordCounts=%s", urlCategoriesEndpoint, url.QueryEscape(customName)), &urlCategory)
	if err != nil {
		return nil, err
	}
	for _, custom := range urlCategory {
		if strings.EqualFold(custom.ID, customName) {
			return &custom, nil
		}
	}
	return nil, fmt.Errorf("no custom url category found with name: %s", customName)
}

func (service *Service) CreateURLCategories(category *URLCategory) (*URLCategory, error) {
	resp, err := service.Client.Create(urlCategoriesEndpoint, *category)
	if err != nil {
		return nil, err
	}

	createdUrlCategory, ok := resp.(*URLCategory)
	if !ok {
		return nil, errors.New("object returned from API was not a url category Pointer")
	}

	service.Client.Logger.Printf("[DEBUG]Returning url category from Create: %v", createdUrlCategory.ID)
	return createdUrlCategory, nil
}

func (service *Service) UpdateURLCategories(categoryID string, category *URLCategory) (*URLCategory, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%s", urlCategoriesEndpoint, categoryID), *category)
	if err != nil {
		return nil, nil, err
	}
	updatedUrlCategory, _ := resp.(*URLCategory)
	service.Client.Logger.Printf("[DEBUG]Returning url category from Update: %s", updatedUrlCategory.ID)
	return updatedUrlCategory, nil, nil
}

func (service *Service) DeleteURLCategories(categoryID string) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%s", urlCategoriesEndpoint, categoryID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (service *Service) GetAll() ([]URLCategory, error) {
	var urlCategories []URLCategory
	err := common.ReadAllPages(service.Client, urlCategoriesEndpoint, &urlCategories)
	return urlCategories, err
}
