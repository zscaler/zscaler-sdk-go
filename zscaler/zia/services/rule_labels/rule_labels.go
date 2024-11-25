package rule_labels

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	ruleLabelsEndpoint = "/zia/api/v1/ruleLabels"
)

type RuleLabels struct {
	// The unique identifier for the rule label.
	ID int `json:"id"`

	// The rule label name.
	Name string `json:"name,omitempty"`

	// The rule label description.
	Description string `json:"description,omitempty"`

	// Timestamp when the rule lable was last modified. This is a read-only field. Ignored by PUT and DELETE requests.
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// The admin that modified the rule label last. This is a read-only field. Ignored by PUT requests.
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// The admin that created the rule label. This is a read-only field. Ignored by PUT requests.
	CreatedBy *common.IDNameExtensions `json:"createdBy,omitempty"`

	// The number of rules that reference the label.
	ReferencedRuleCount int `json:"referencedRuleCount,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, ruleLabelID int) (*RuleLabels, error) {
	var ruleLabel RuleLabels
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", ruleLabelsEndpoint, ruleLabelID), &ruleLabel)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning rule label from Get: %d", ruleLabel.ID)
	return &ruleLabel, nil
}

func GetRuleLabelByName(ctx context.Context, service *zscaler.Service, labelName string) (*RuleLabels, error) {
	var ruleLabels []RuleLabels
	err := common.ReadAllPages(ctx, service.Client, ruleLabelsEndpoint, &ruleLabels)
	if err != nil {
		return nil, err
	}
	for _, ruleLabel := range ruleLabels {
		if strings.EqualFold(ruleLabel.Name, labelName) {
			return &ruleLabel, nil
		}
	}
	return nil, fmt.Errorf("no rule label found with name: %s", labelName)
}

func Create(ctx context.Context, service *zscaler.Service, ruleLabelID *RuleLabels) (*RuleLabels, *http.Response, error) {
	resp, err := service.Client.Create(ctx, ruleLabelsEndpoint, *ruleLabelID)
	if err != nil {
		return nil, nil, err
	}

	createdRuleLabel, ok := resp.(*RuleLabels)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a rule label pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new rule label from create: %d", createdRuleLabel.ID)
	return createdRuleLabel, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleLabelID int, ruleLabels *RuleLabels) (*RuleLabels, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", ruleLabelsEndpoint, ruleLabelID), *ruleLabels)
	if err != nil {
		return nil, nil, err
	}
	updatedRuleLabel, _ := resp.(*RuleLabels)

	service.Client.GetLogger().Printf("[DEBUG]returning updates rule label from update: %d", updatedRuleLabel.ID)
	return updatedRuleLabel, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleLabelID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", ruleLabelsEndpoint, ruleLabelID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]RuleLabels, error) {
	var ruleLabels []RuleLabels
	err := common.ReadAllPages(ctx, service.Client, ruleLabelsEndpoint, &ruleLabels)
	return ruleLabels, err
}
