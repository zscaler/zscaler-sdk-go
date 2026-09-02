package users

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ziam/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ziam/services/groups"
)

const (
	usersEndpoint = "/ziam/admin/api/v1/users"
)

type Users struct {
	ID             string `json:"id,omitempty"`
	Source         string `json:"source,omitempty"`
	LoginName      string `json:"loginName,omitempty"`
	DisplayName    string `json:"displayName,omitempty"`
	FirstName      string `json:"firstName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	PrimaryEmail   string `json:"primaryEmail,omitempty"`
	SecondaryEmail string `json:"secondaryEmail,omitempty"`

	// Status is a pointer so that disabling a user is expressible.
	//
	// As a plain bool with `omitempty`, false was indistinguishable from unset
	// and was dropped from the request body, which made it possible to enable a
	// user but never to disable one. A nil pointer still omits the field, so
	// callers that do not care about status are unaffected.
	Status *bool `json:"status,omitempty"`

	Department      *common.IDNameDisplayName `json:"department,omitempty"`
	IDP             *common.IDNameDisplayName `json:"idp,omitempty"`
	CustomAttrsInfo map[string]interface{}    `json:"customAttrsInfo,omitempty"`
}

type UsersResponse = common.PaginationResponse[Users]

func GetUser(ctx context.Context, service *zscaler.Service, userID string) (*Users, error) {
	var user Users
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%s", usersEndpoint, userID), &user)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning user from Get: %s", user.ID)
	return &user, nil
}

// GetAll retrieves all users with optional pagination and filtering parameters
func GetAll(ctx context.Context, service *zscaler.Service, queryParams *common.PaginationQueryParams) ([]Users, error) {
	return common.ReadAllPagesWithPagination[Users](ctx, service.Client, usersEndpoint, queryParams)
}

// GetByName retrieves users by searching through paginated data for the specified name
func GetByName(ctx context.Context, service *zscaler.Service, name string) ([]Users, error) {
	var allUsers []Users
	var currentOffset int
	pageSize := 100 // Use a reasonable page size for searching

	for {
		// Create query params for current page
		queryParams := common.NewPaginationQueryParams(pageSize)
		queryParams.WithOffset(currentOffset)

		// Get current page
		pageResponse, err := common.ReadPageWithPagination[Users](ctx, service.Client, usersEndpoint, &queryParams)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page at offset %d: %w", currentOffset, err)
		}

		// Search through records in this page
		for _, user := range pageResponse.Records {
			if strings.Contains(strings.ToLower(user.DisplayName), strings.ToLower(name)) {
				allUsers = append(allUsers, user)
			}
		}

		// Check if we've reached the end
		if len(pageResponse.Records) < pageSize || pageResponse.NextLink == "" {
			break
		}

		// Move to next page
		currentOffset += len(pageResponse.Records)
	}

	return allUsers, nil
}

func Create(ctx context.Context, service *zscaler.Service, user *Users) (*Users, *http.Response, error) {
	resp, err := service.Client.Create(ctx, usersEndpoint, *user)
	if err != nil {
		return nil, nil, err
	}

	createdUser, ok := resp.(*Users)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a user pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new user from create: %s", createdUser.ID)
	return createdUser, nil, nil
}

// Update replaces a user via PUT.
//
// The response is nil-checked before use: the shared client returns a nil
// object for an empty body, so a `204 No Content` answer to the PUT would
// otherwise dereference nil and panic the calling process. A nil return with a
// nil error means the update succeeded but the endpoint reported no body — the
// caller should re-read the user if it needs the current representation.
func Update(ctx context.Context, service *zscaler.Service, userID string, user *Users) (*Users, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%s", usersEndpoint, userID), *user)
	if err != nil {
		return nil, nil, err
	}

	updatedUser, ok := resp.(*Users)
	if !ok || updatedUser == nil {
		service.Client.GetLogger().Printf("[DEBUG]update of user %s returned no body", userID)
		return nil, nil, nil
	}

	service.Client.GetLogger().Printf("[DEBUG]returning updated user from update: %s", updatedUser.ID)
	return updatedUser, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, userID string) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%s", usersEndpoint, userID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetGroupsByUser retrieves every group a user belongs to, across all pages.
//
// This previously read only the first page while being documented as returning
// all of them, which silently truncated the answer for any user in more than
// `limit` groups — a caller reconciling state against the result would read the
// truncation as the user having been removed from the missing groups.
func GetGroupsByUser(ctx context.Context, service *zscaler.Service, userID string, queryParams *common.PaginationQueryParams) ([]groups.Groups, error) {
	groupsEndpoint := fmt.Sprintf("%s/%s/groups", usersEndpoint, userID)
	return common.ReadAllPagesWithPagination[groups.Groups](ctx, service.Client, groupsEndpoint, queryParams)
}

// GetGroupsByUserPage retrieves a single page of a user's group memberships,
// along with the pagination metadata. Use GetGroupsByUser unless you need to
// page manually.
func GetGroupsByUserPage(ctx context.Context, service *zscaler.Service, userID string, queryParams *common.PaginationQueryParams) (*common.PaginationResponse[groups.Groups], error) {
	groupsEndpoint := fmt.Sprintf("%s/%s/groups", usersEndpoint, userID)
	return common.ReadPageWithPagination[groups.Groups](ctx, service.Client, groupsEndpoint, queryParams)
}

// =============================================================================
// User operations
// =============================================================================
//
// These four endpoints are actions rather than objects: they have a POST or PUT
// and nothing else. There is no GET, so nothing they set can be read back, and
// no DELETE, so nothing they do can be undone through the API.
//
// They deliberately bypass service.Client.Create and service.Client.UpdateWithPut
// in favour of ExecuteRequest. Those two helpers unmarshal the response body
// into a new value of the *request* struct's type, and every one of these
// endpoints documents its 200 body as a bare JSON string. Feeding `"Success"`
// into a *SkipMFARequest fails, so the helpers would report an error for a call
// that in fact succeeded. ExecuteRequest hands back the raw bytes and leaves the
// interpretation here.

// SkipMFARequest is the body of SetSkipMFA.
//
// Timestamp is not tagged `omitempty`: it is required, and a zero epoch is a
// value the caller may legitimately want to send.
type SkipMFARequest struct {
	Timestamp int64 `json:"timestamp"`
}

// UpdatePasswordRequest is the body of UpdatePassword.
//
// ResetPwdOnLogin is not tagged `omitempty`, so an explicit false is always
// transmitted. The alternative repeats the mistake that made Users.Status a
// pointer: with `omitempty` a caller could require a reset on next login but
// never stop requiring one.
type UpdatePasswordRequest struct {
	Password        string `json:"password"`
	ResetPwdOnLogin bool   `json:"resetPwdOnLogin"`
}

// SetSkipMFA sets the flag that lets a user skip Multi-Factor Authentication.
//
// The timestamp is documented as an int32 holding epoch time, and records when
// the skip was set. It is accepted here as an int64 so that callers are not
// forced through a lossy conversion; a tenant that enforces the documented
// int32 will reject an out-of-range value with an API error rather than
// silently truncating.
//
// There is no companion endpoint to clear the flag.
func SetSkipMFA(ctx context.Context, service *zscaler.Service, userID string, timestamp int64) (*http.Response, error) {
	return executeUserAction(ctx, service, http.MethodPost,
		fmt.Sprintf("%s/%s/setskipmfa", usersEndpoint, userID),
		SkipMFARequest{Timestamp: timestamp},
	)
}

// ResetPassword initiates a password reset for a user, which in practice means
// ZIdentity sends the user a reset invitation.
//
// The endpoint takes no request body. It is an action with no state: calling it
// twice sends two invitations, and there is nothing to read back afterwards.
func ResetPassword(ctx context.Context, service *zscaler.Service, userID string) (*http.Response, error) {
	return executeUserAction(ctx, service, http.MethodPost,
		fmt.Sprintf("%s/%s/resetpassword", usersEndpoint, userID),
		nil,
	)
}

// UpdatePassword sets a user's password, optionally requiring them to change it
// at their next login.
//
// Note this is a PUT while the other two actions are POSTs.
func UpdatePassword(ctx context.Context, service *zscaler.Service, userID string, payload *UpdatePasswordRequest) (*http.Response, error) {
	if payload == nil {
		return nil, errors.New("tried to update a password with a nil payload")
	}

	return executeUserAction(ctx, service, http.MethodPut,
		fmt.Sprintf("%s/%s/updatepassword", usersEndpoint, userID),
		*payload,
	)
}

// executeUserAction issues one of the action requests above.
//
// A nil body sends no request body at all, which is what /resetpassword
// expects. The response body is logged rather than returned: all three
// endpoints document it as an opaque status string, so there is nothing in it a
// caller can act on beyond the absence of an error.
func executeUserAction(ctx context.Context, service *zscaler.Service, method, endpoint string, body interface{}) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}

	respBody, resp, _, err := service.Client.ExecuteRequest(ctx, method, endpoint, reader, nil, "application/json")
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]%s %s returned: %s", method, endpoint, string(respBody))
	return resp, nil
}
