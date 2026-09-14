package middleware

import "strings"

// IsXWebView reports whether the user agent identifies the X app's web view.
func IsXWebView(userAgent string) bool {
	userAgent = strings.ToLower(userAgent)

	return strings.Contains(userAgent, "twitter") ||
		strings.Contains(userAgent, "xandroid") ||
		strings.Contains(userAgent, "xwebview")
}
