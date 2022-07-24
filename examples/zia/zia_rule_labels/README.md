```go
package main

import (
	"log"

	"github.com/zscaler/zscaler-sdk-go/zia"
	"github.com/zscaler/zscaler-sdk-go/zia/services/rule_labels"
)

func main() {
	cli, err := zia.NewClient("username@acme.com", "pwd", "apiKey", "zscalerthree", "userAgent")
	if err != nil {
		log.Printf("[ERROR] creating client failed: %v\n", err)
		return
	}
	ruleLabelsService := rule_labels.New(cli)
	ruleLabels := rule_labels.RuleLabels{
		Name:              "Example rule labels",
		Description:       "Example rule labels",
	}
	// Create new rule labels
	createRuleLabels, err := ruleLabelsService.Create(&ruleLabels)
	if err != nil {
		log.Printf("[ERROR] creating rule labels failed: %v\n", err)
		return
	}
	// Update rule labels
	createRuleLabels.Description = "New description"
	_, _, err = ruleLabelsService.Update(createRuleLabels.ID, createRuleLabels)
	if err != nil {
		log.Printf("[ERROR] updating rule labels failed: %v\n", err)
		return
	}
	// Delete rule labels
	_, err = ruleLabelsService.Delete(createRuleLabels.ID)
	if err != nil {
		log.Printf("[ERROR] deleting rule labels failed: %v\n", err)
		return
	}
}

```
