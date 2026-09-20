package middleware

import "strings"

// isLinkPreviewBot reports whether the user agent identifies a crawler that
// fetches pages to build link previews (e.g. X/Twitter's card bot). These
// must always see the real page with OGP tags, never the webview guide.
func isLinkPreviewBot(userAgent string) bool {
	return strings.Contains(userAgent, "twitterbot") ||
		strings.Contains(userAgent, "facebookexternalhit") ||
		strings.Contains(userAgent, "discordbot") ||
		strings.Contains(userAgent, "slackbot") ||
		strings.Contains(userAgent, "linkedinbot") ||
		strings.Contains(userAgent, "whatsapp")
}

// IsXWebView reports whether the user agent identifies the X app's web view.
func IsXWebView(userAgent string) bool {
	userAgent = strings.ToLower(userAgent)

	if isLinkPreviewBot(userAgent) {
		return false
	}

	return strings.Contains(userAgent, "twitter") ||
		strings.Contains(userAgent, "xandroid") ||
		strings.Contains(userAgent, "xwebview")
}
