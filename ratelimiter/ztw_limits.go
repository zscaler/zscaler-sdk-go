package ratelimiter

// Cloud & Branch Connector (ZTW) rate limits. The Cloud & Branch Connector API
// documents the same weights as ZIA (GET 2/sec and 1000/hr, POST and PUT 1/sec
// and 400/hr) and runs on the same platform, so it uses the ZIA limits: the
// same defaults per method and resource URI, and the same stricter resources
// where a ZTW resource shares the ZIA path (for example GET /dnsGateways).

// ZTWPathPrefix is the prefix of Cloud & Branch Connector paths on the OneAPI
// gateway.
const ZTWPathPrefix = "/ztw/api/v1"

// NewZTWEndpointRateLimiter returns a limiter that applies the ZIA limits to
// Cloud & Branch Connector requests, per method and resource URI.
func NewZTWEndpointRateLimiter() *EndpointRateLimiter {
	return NewEndpointRateLimiter(ZTWPathPrefix, ZIADefaultLimits(), ZIAEndpointRules())
}
