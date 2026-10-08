---
page_title: "Troubleshooting Guide"
---

# How to troubleshoot your problem

## Collect debug logs

Enable SDK logging before reproducing the problem:

```sh
export ZSCALER_SDK_LOG=true
export ZSCALER_SDK_VERBOSE=true
```

The log records every request and the full response, including the `x-oneapi-request-id` and `x-oneapi-host` response headers. Include those headers and the UTC time of the failing request when opening a support case; they are required to trace a request on the Zscaler side.

Credentials in request and response headers (`Authorization`, `Cookie`, `Set-Cookie`, `JSessionID`, `auth-token`) are masked in the log from v3.8.53. Earlier versions log the OAuth access token in clear text.

~> **IMPORTANT:** Sanitize the log before sharing it. It still contains tenant identifiers and request and response data. Remove or mask these values before attaching the log to a public GitHub issue, or provide it through a formal Zscaler support case.

## How errors are reported

API errors are returned as `*errorx.ErrorResponse`. Its `Parsed` field holds the `code`, `message`, `id`, `reason` and `exception` values from the API's JSON error body, plus the request URL and HTTP status. When the body contains none of these fields, `message` holds the body exactly as the API returned it (v3.8.53+). An empty `message` therefore means the API returned an empty body.

Use `errorx.AsErrorResponse(err)` to extract it, and helpers such as `IsObjectNotFound()` and `IsLimitExceeded()` to classify it.

## Automatic retries

The OneAPI client retries the following on its own; callers should not add their own retry loop.

| Response | Behaviour |
|---|---|
| `429`, `503` | Waits at least the `Retry-After` interval, growing the wait on consecutive retries. A `Retry-After` longer than 5 minutes fails immediately. |
| `409`, `412` edit lock / org barrier | Retried with backoff. |
| `5xx` | Retried with backoff, except when the body carries an API error code (for example `500 UNEXPECTED_ERROR` for invalid input). `502`, `503` and `504` are always retried; `501` never is. |
| `401` or `403` indicating an ended or invalid session (`SESSION_NOT_VALID`, `Resource Access Blocked`, `Session already invalidated`) | Obtains a new OAuth token and retries, up to `MaxSessionNotValidRetries` times (default 3, env `ZSCALER_CLIENT_MAX_SESSION_NOT_VALID_RETRIES`). `403` is handled this way from v3.8.53. |
| Any other `4xx` | Not retried. |

## 403 `Resource Access Blocked`

Two common causes:

- **Every request for a given resource type fails, from the start:** the role assigned to the API client lacks permission or functional scope for it, or the tenant is missing a required subscription. A new token does not help; the request fails once the session retries are used.
- **Requests start failing several minutes into a long-running process, and a new token restores access:** the API session ended. ZIA ends API-initiated sessions after the API Session Timeout in Advanced Settings (5 to 20 minutes, default 5) even though the OAuth token has not expired. The client obtains a new token and retries automatically.

Other `403` responses, such as a missing permission, a missing subscription, a write during a scheduled maintenance window, or `LIMIT_EXCEEDED`, are returned immediately.

## Intermittent 401 with an empty `code` and `message`

The API returned a `401` whose body carries no error details. The client does not retry it, because the cause cannot be determined from the response. Capture a debug log of a failing request and open a support case with the `x-oneapi-request-id` and `x-oneapi-host` headers.

If many processes use the same API client against the same tenant at the same time, reduce that concurrency: each process paces its own requests, so together they can exceed the tenant's API limits.
