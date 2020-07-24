package matomo

type Options struct {
	// The URL of your matomo installation (with our without /piwik.php)
	PiwikURL string
	// Ignore the Do not Track header that is sent by the browser. This is not recommended
	IgnoreDoNotTrack bool
	// The ID of the website in piwik
	WebsiteID string
	// The piwik API's access token
	Token string
}
