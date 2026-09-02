package apiclients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ziam/services/common"
)

const (
	apiClientsEndpoint = "/ziam/admin/api/v1/api-clients"
)

// AuthType is the method an API client authenticates with.
//
// The three values are mutually exclusive and each one determines which of
// ClientAuthentication's other fields the API expects:
//
//   - AuthTypeSecret     — the client authenticates with a secret, managed
//     through AddSecret / GetSecrets / DeleteSecret.
//   - AuthTypePubKeyCert — the client authenticates with a certificate, supplied
//     in ClientCertificates or PublicKeys.
//   - AuthTypeJWKS       — the client authenticates with keys published at
//     ClientJWKsURL.
type AuthType string

const (
	AuthTypeSecret     AuthType = "SECRET"
	AuthTypePubKeyCert AuthType = "PUBKEYCERT"
	AuthTypeJWKS       AuthType = "JWKS"
)

// APIClients is an OAuth2 API client.
//
// As elsewhere in this package set, one struct serves both the request and the
// response. The response carries three fields the request does not — Name and
// DefaultApi on ClientResource, and Name on SelectedScope — and all three are
// `omitempty`, so sending a client read back from the API costs nothing but a
// few extra keys the server ignores. Omitting them instead would silently drop
// data the endpoint does return.
type APIClients struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	// Status is a pointer so that disabling a client is expressible.
	//
	// As a plain bool with `omitempty`, false would be indistinguishable from
	// unset and dropped from the request body, making it possible to enable a
	// client but never to disable one. This is the same defect that made
	// users.Users.Status a pointer. A nil pointer still omits the field, so
	// callers that do not care about status are unaffected.
	Status *bool `json:"status,omitempty"`

	// AccessTokenLifeTime is how long an access token issued to this client
	// remains valid, and is required on both create and update. It carries no
	// `omitempty`: a PUT is a full replace, so a caller that leaves it at zero
	// should get the API's range error rather than have the field silently
	// vanish and the server apply a default.
	//
	// The documented range is 60 to 86400. The unit is contradictory in the
	// official documentation — the create and response schemas say seconds
	// while the update schema says minutes, and both give "60 for 60 minutes"
	// as the example. Seconds is the reading consistent with the 86400 ceiling
	// (24 hours), but this has not been confirmed against a tenant, so treat
	// the unit as unverified rather than relying on it.
	AccessTokenLifeTime int32 `json:"accessTokenLifeTime"`

	// ClientAuthentication is required. It is a pointer because `omitempty` has
	// no effect on a struct value: as a value type, an omitted
	// `clientAuthentication` would unmarshal to an empty struct rather than
	// nil, which is indistinguishable from one whose fields are genuinely
	// blank.
	ClientAuthentication *ClientAuthentication `json:"clientAuthentication,omitempty"`

	// ClientResources is required, and lists the resource servers the client may
	// reach along with the scopes selected for each. Resolve the resource and
	// scope ids with the resource_servers package.
	ClientResources []ClientResource `json:"clientResources,omitempty"`
}

// ClientAuthentication holds how a client proves its identity. Which fields the
// API expects depends on AuthType.
type ClientAuthentication struct {
	AuthType AuthType `json:"authType,omitempty"`

	// ClientJWKsURL is where the client publishes its JSON Web Key Set. API
	// field: clientJWKsUrl.
	ClientJWKsURL string `json:"clientJWKsUrl,omitempty"`

	PublicKeys         []PublicKey         `json:"publicKeys,omitempty"`
	ClientCertificates []ClientCertificate `json:"clientCertificates,omitempty"`
}

// PublicKey is one named public key.
type PublicKey struct {
	KeyName  string `json:"keyName,omitempty"`
	KeyValue string `json:"keyValue,omitempty"`
}

// ClientCertificate is one client certificate, carried as a string.
type ClientCertificate struct {
	CertContent string `json:"certContent,omitempty"`
}

// ClientResource is a resource server the client may reach, and the scopes
// selected on it.
//
// Only ID and SelectedScopes are accepted on write. Name and DefaultApi are
// reported on read.
type ClientResource struct {
	ID             string          `json:"id,omitempty"`
	Name           string          `json:"name,omitempty"`
	DefaultApi     bool            `json:"defaultApi,omitempty"`
	SelectedScopes []SelectedScope `json:"selectedScopes,omitempty"`
}

// SelectedScope is one scope selected on a resource. Only ID is accepted on
// write; Name is reported on read.
//
// ID is the composite id that resource_servers reports verbatim, and it is
// passed through unchanged. Its internal shape is "<serviceId>:<zpaScopeId>:<id>"
// with the middle segment empty for every service except ZPA — so ZIA reports
// "hhlm44raf07ps::hplm45bceg7r0" and ZPA "hhlm44rae07ib:mplm44rqi07jb:hplm44rqvg7n5"
// — but that is an implementation detail of the id, not a structure to take
// apart. This endpoint accepts and returns the joined form on both read and
// write; splitting it into separate tserviceId / zpaScopeId fields is rejected
// with "400 Invalid Input for clientResources". Confirmed against a tenant.
type SelectedScope struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// APIClientSecret is a secret belonging to an API client.
//
// Value is only ever returned in full by AddSecret. GetSecrets masks it, so a
// secret that is not recorded at the moment it is created cannot be recovered.
//
// Both timestamps are UNIX epochs in SECONDS. The official API reference states
// this explicitly for expiresAt, in the create request body and in the create
// and list response schemas alike.
//
// The ZIdentity admin console reports them in milliseconds — a secret created at
// 2026-09-01T21:21:06Z shows up there as 1788297666033 — but the console is not
// this API, and its payloads have already proven a poor guide to it: it also
// sends a selected scope's id split across three fields, which this endpoint
// rejects with 400. Do not "correct" these fields to milliseconds on the
// strength of a console capture.
//
// The two are int64 rather than int, and must stay that way. The API documents
// both as int64, `int` is 32 bits on the 386 and arm builds the providers ship,
// and any expiry past January 2038 exceeds MaxInt32 under those builds and
// fails to unmarshal with "number out of range". A five-year secret issued
// today already crosses that line. The sibling users.SkipMFARequest.Timestamp
// is int64 for the same reason.
type APIClientSecret struct {
	ID        string `json:"id,omitempty"`
	Value     string `json:"value,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	ExpiresAt int64  `json:"expiresAt,omitempty"`
}

// AddSecretRequest is the body of AddSecret.
//
// ExpiresAt is required and carries no `omitempty`, matching
// users.SkipMFARequest.Timestamp: a zero epoch is a value the caller may
// legitimately send, and dropping a required field turns a caller bug into an
// opaque server-side default.
type AddSecretRequest struct {
	ExpiresAt int64 `json:"expiresAt"`
}

type APIClientsResponse = common.PaginationResponse[APIClients]

func Get(ctx context.Context, service *zscaler.Service, clientID string) (*APIClients, error) {
	var apiClient APIClients
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%s", apiClientsEndpoint, clientID), &apiClient)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning api client from Get: %s", apiClient.ID)
	return &apiClient, nil
}

// GetAll retrieves all API clients with optional pagination and filtering parameters
func GetAll(ctx context.Context, service *zscaler.Service, queryParams *common.PaginationQueryParams) ([]APIClients, error) {
	return common.ReadAllPagesWithPagination[APIClients](ctx, service.Client, apiClientsEndpoint, queryParams)
}

// GetPage retrieves a single page of API clients along with the pagination
// metadata — results_total, next_link, and prev_link. Use GetAll unless you
// need to page manually or want the record count.
func GetPage(ctx context.Context, service *zscaler.Service, queryParams *common.PaginationQueryParams) (*APIClientsResponse, error) {
	return common.ReadPageWithPagination[APIClients](ctx, service.Client, apiClientsEndpoint, queryParams)
}

// GetByName retrieves API clients by searching through paginated data for the specified name
//
// This is a case-insensitive substring match applied client-side, so a search
// for "prod" also returns "production-api" and "reprod". Callers that need one
// specific client must narrow the result to an exact match themselves.
//
// The endpoint does accept a server-side `name[like]` filter, reachable through
// common.PaginationQueryParams.WithNameFilter on GetAll. It is not used here
// because its matching semantics are undocumented, and a filter that turned out
// to be prefix-only or case-sensitive would drop matches this function returns.
func GetByName(ctx context.Context, service *zscaler.Service, name string) ([]APIClients, error) {
	var allAPIClients []APIClients
	var currentOffset int
	pageSize := 100 // Use a reasonable page size for searching

	for {
		// Create query params for current page
		queryParams := common.NewPaginationQueryParams(pageSize)
		queryParams.WithOffset(currentOffset)

		// Get current page
		pageResponse, err := common.ReadPageWithPagination[APIClients](ctx, service.Client, apiClientsEndpoint, &queryParams)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page at offset %d: %w", currentOffset, err)
		}

		// Search through records in this page
		for _, apiClient := range pageResponse.Records {
			if strings.Contains(strings.ToLower(apiClient.Name), strings.ToLower(name)) {
				allAPIClients = append(allAPIClients, apiClient)
			}
		}

		// Check if we've reached the end
		if len(pageResponse.Records) < pageSize || pageResponse.NextLink == "" {
			break
		}

		// Move to next page
		currentOffset += len(pageResponse.Records)
	}

	return allAPIClients, nil
}

func Create(ctx context.Context, service *zscaler.Service, apiClient *APIClients) (*APIClients, *http.Response, error) {
	if apiClient == nil {
		return nil, nil, errors.New("tried to create an api client with a nil payload")
	}

	resp, err := service.Client.Create(ctx, apiClientsEndpoint, *apiClient)
	if err != nil {
		return nil, nil, err
	}

	createdAPIClient, ok := resp.(*APIClients)
	if !ok {
		return nil, nil, errors.New("object returned from api was not an api client pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new api client from create: %s", createdAPIClient.ID)
	return createdAPIClient, nil, nil
}

// Update replaces an API client via PUT.
//
// The response is nil-checked before use: the shared client returns a nil
// object for an empty body, so a `204 No Content` answer to the PUT would
// otherwise dereference nil and panic the calling process. A nil return with a
// nil error means the update succeeded but the endpoint reported no body — the
// caller should re-read the client if it needs the current representation.
//
// A PUT is a full replace, so build the payload by reading the live client and
// overlaying the intended changes. Sending only the changed fields drops
// everything omitted.
func Update(ctx context.Context, service *zscaler.Service, clientID string, apiClient *APIClients) (*APIClients, *http.Response, error) {
	if apiClient == nil {
		return nil, nil, errors.New("tried to update an api client with a nil payload")
	}

	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%s", apiClientsEndpoint, clientID), *apiClient)
	if err != nil {
		return nil, nil, err
	}

	updatedAPIClient, ok := resp.(*APIClients)
	if !ok || updatedAPIClient == nil {
		service.Client.GetLogger().Printf("[DEBUG]update of api client %s returned no body", clientID)
		return nil, nil, nil
	}

	service.Client.GetLogger().Printf("[DEBUG]returning updated api client from update: %s", updatedAPIClient.ID)
	return updatedAPIClient, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, clientID string) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%s", apiClientsEndpoint, clientID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// =============================================================================
// Client secrets
// =============================================================================

// GetSecrets retrieves the secrets associated with an API client.
//
// Two things to note. The endpoint answers with a bare JSON array rather than
// the paginated envelope the rest of this package uses, so there is no page
// variant. And the secret values are masked, so this reports which secrets
// exist and when they expire, not what they are.
//
// Only meaningful when the client's AuthType is AuthTypeSecret.
func GetSecrets(ctx context.Context, service *zscaler.Service, clientID string) ([]APIClientSecret, error) {
	var secrets []APIClientSecret
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%s/secrets", apiClientsEndpoint, clientID), &secrets)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning %d secret(s) for api client %s", len(secrets), clientID)
	return secrets, nil
}

// AddSecret creates a secret and associates it with an API client, returning it
// with its value populated.
//
// expiresAt is the expiry as a UNIX epoch timestamp in seconds, and is
// required.
//
// This bypasses service.Client.Create in favour of ExecuteRequest, and the
// reason is not stylistic. Create unmarshals the response body into a new value
// of the *request* struct's type, which here is AddSecretRequest — a struct
// with nothing but ExpiresAt. The response carries the id, the creation time,
// and the secret value, so routing it through Create would quietly discard the
// value: the one field the call exists to obtain, and the one field no
// subsequent request can recover, since GetSecrets masks it.
//
// Only meaningful when the client's AuthType is AuthTypeSecret.
func AddSecret(ctx context.Context, service *zscaler.Service, clientID string, expiresAt int64) (*APIClientSecret, *http.Response, error) {
	endpoint := fmt.Sprintf("%s/%s/secrets", apiClientsEndpoint, clientID)

	data, err := json.Marshal(AddSecretRequest{ExpiresAt: expiresAt})
	if err != nil {
		return nil, nil, err
	}

	respBody, resp, _, err := service.Client.ExecuteRequest(ctx, http.MethodPost, endpoint, bytes.NewReader(data), nil, "application/json")
	if err != nil {
		return nil, nil, err
	}

	// A 201 carries the new secret. An empty body is treated as success with
	// nothing to report rather than an error, so that a tenant answering 204
	// does not surface as a failed call — but the caller is warned, because
	// without the body the secret value is unrecoverable.
	if len(respBody) == 0 {
		service.Client.GetLogger().Printf("[WARN]creation of a secret for api client %s returned no body; the secret value cannot be recovered", clientID)
		return nil, resp, nil
	}

	var secret APIClientSecret
	if err := json.Unmarshal(respBody, &secret); err != nil {
		return nil, nil, fmt.Errorf("failed to parse the secret created for api client %s: %w", clientID, err)
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new secret from create: %s", secret.ID)
	return &secret, resp, nil
}

// DeleteSecret removes one secret from an API client. The secret cannot be
// recovered afterwards.
func DeleteSecret(ctx context.Context, service *zscaler.Service, clientID, secretID string) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%s/secrets/%s", apiClientsEndpoint, clientID, secretID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
