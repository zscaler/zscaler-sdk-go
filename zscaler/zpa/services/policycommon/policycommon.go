// Package policycommon holds types shared across the ZPA policy packages
// (policysetcontroller, policysetcontrollerv2, policy_group, policy_group_rule,
// policy_group_set, ...).
//
// These types cannot live in zpa/services/common because they reference
// service packages such as applicationsegment, and applicationsegment already
// imports common — placing them in common would create an import cycle. This
// package sits above both common and applicationsegment, so every policy
// package can import it without cycles.
package policycommon

import "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"

type DesktopPolicyMappings struct {
	AppSegments  []applicationsegment.ApplicationSegmentResource `json:"appSegments,omitempty"`
	ID           string                                          `json:"id"`
	CreationTime string                                          `json:"creationTime,omitempty"`
	ModifiedBy   string                                          `json:"modifiedBy,omitempty"`
	ModifiedTime string                                          `json:"modifiedTime,omitempty"`
	ImageID      string                                          `json:"imageId,omitempty"`
	ImageName    string                                          `json:"imageName,omitempty"`
}
