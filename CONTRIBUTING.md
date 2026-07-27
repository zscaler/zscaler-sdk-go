# Contributing

## This repository does not accept external contributions

The Zscaler SDK for Go is developed and maintained internally by Zscaler against
internal API specifications and release processes. **We do not accept pull
requests.**

This applies to every kind of change, including bug fixes, new services and
endpoints, added struct fields, documentation corrections and dependency
updates. Pull requests opened against this repository will be closed without
review.

This is not a judgement on the quality of any proposed change. The SDK's service
definitions, versioning and release pipeline are driven by internal sources that
a pull request cannot participate in, so a merged external patch would be
overwritten by the next generated release.

## How to request a fix or a change

**Open a GitHub issue.** That is the only supported channel, and it is the
channel the maintainers actually work from.

| I want to | Open |
|---|---|
| Report something broken or behaving contrary to the docs | [🐛 Bug report](https://github.com/zscaler/zscaler-sdk-go/issues/new?template=bug.yml) |
| Request a new service, endpoint, field or capability | [🚀 Feature request](https://github.com/zscaler/zscaler-sdk-go/issues/new?template=feature_request.yml) |

Issues are triaged by the maintainers and, when accepted, implemented and
shipped in a subsequent release. You can track the outcome in
[CHANGELOG.md](CHANGELOG.md) and the
[release notes](docs/guides/release-notes.md).

### Writing an issue we can act on

The faster we can reproduce a problem, the faster it gets fixed. Please include:

- The **SDK version** you are on (`github.com/zscaler/zscaler-sdk-go/v3 vX.Y.Z`
  from your `go.mod`) and your Go version.
- The **cloud and service** involved — ZIA, ZPA, ZCC, ZDX, ZTW, ZID or ZWA —
  and whether you are using the OneAPI client or a legacy client.
- A **minimal reproduction**: the SDK call you made and the arguments, reduced
  to the smallest form that still shows the problem.
- The **full error**, including the HTTP status and the API's `code` and
  `message` where present.
- What you expected to happen instead.

Redact credentials, tokens, customer IDs and any tenant-identifying data before
posting. Issues are public.

### Security issues

Do not report suspected security vulnerabilities in a public issue or pull
request. Contact Zscaler through the channels described in
[SUPPORT.md](SUPPORT.md).

## Support scope

This project is released under the as-is, best-effort policy described in
[SUPPORT.md](SUPPORT.md). Please read it before opening an issue, so the level
of support you can expect is clear.

## Forking for your own use

The license permits you to fork this repository and modify your copy freely. You
are welcome to do so — just note that we will not merge those changes back, and
a fork will not receive upstream releases automatically. If your change would
benefit other users, open a feature request describing the need rather than the
patch, and we will evaluate it for the SDK itself.
