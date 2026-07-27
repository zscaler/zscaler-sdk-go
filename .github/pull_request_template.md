<!--
  ┌──────────────────────────────────────────────────────────────────────┐
  │  STOP — this repository does not accept pull requests.               │
  │  Please close this PR and open a GitHub issue instead.               │
  └──────────────────────────────────────────────────────────────────────┘
-->

# ⛔ This repository does not accept pull requests

**Please do not submit this pull request.** It will be closed without review.

The Zscaler SDK for Go is developed and maintained internally by Zscaler against
internal API specifications and release processes. Externally merged changes
would be overwritten by the next generated release, so we cannot accept them —
including bug fixes, new services and endpoints, added fields, documentation
corrections and dependency updates.

This is not a judgement on your change. We do want to hear about it.

## ✅ What to do instead

Open a GitHub issue — the only supported channel, and the one maintainers work
from:

- 🐛 **[Report a bug](https://github.com/zscaler/zscaler-sdk-go/issues/new?template=bug.yml)** — something is broken or behaves contrary to the documentation
- 🚀 **[Request a feature](https://github.com/zscaler/zscaler-sdk-go/issues/new?template=feature_request.yml)** — a new service, endpoint, field or capability

Describe the **problem or the need**, not the patch. Accepted issues are
implemented by the maintainers and shipped in a subsequent release, which you
can track in [CHANGELOG.md](https://github.com/zscaler/zscaler-sdk-go/blob/master/CHANGELOG.md) and the
[release notes](https://github.com/zscaler/zscaler-sdk-go/blob/master/docs/guides/release-notes.md).

To help us reproduce quickly, include your SDK and Go versions, the cloud and
service involved (ZIA, ZPA, ZCC, ZDX, ZTW, ZID, ZWA) and whether you use the
OneAPI or a legacy client, a minimal reproduction, and the full error with the
API's HTTP status, `code` and `message`. **Redact credentials, tokens, customer
IDs and tenant-identifying data — issues are public.**

Found a suspected security vulnerability? Do not open a public issue or PR;
follow [SUPPORT.md](https://github.com/zscaler/zscaler-sdk-go/blob/master/SUPPORT.md) instead.

See [CONTRIBUTING.md](https://github.com/zscaler/zscaler-sdk-go/blob/master/CONTRIBUTING.md) for the full policy and
[SUPPORT.md](https://github.com/zscaler/zscaler-sdk-go/blob/master/SUPPORT.md) for the support scope.

---

<sub>Zscaler maintainers only: this notice replaces the previous PR template. Delete it and describe your change as usual.</sub>
