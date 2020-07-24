package matomo

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	/**
	 * Ecommerce item page view tracking stores item's metadata in these Custom Variables slots.
	 */
	CVAR_INDEX_ECOMMERCE_ITEM_PRICE    = 2
	CVAR_INDEX_ECOMMERCE_ITEM_SKU      = 3
	CVAR_INDEX_ECOMMERCE_ITEM_NAME     = 4
	CVAR_INDEX_ECOMMERCE_ITEM_CATEGORY = 5
	/**
	 * Defines how many categories can be used max when calling addEcommerceItem().
	 * @var int
	 */
	MAX_NUM_ECOMMERCE_ITEM_CATEGORIES = 5
	DEFAULT_COOKIE_PATH               = '/'
)

type Matomo struct {
	Options *Options `json:"options"`
}

type TrackingParams struct {
	// custom variables that will be assigned to the action (cvar)
	ActionCVar map[string]string

	// custom variables that will be assigned to the visitor (_cvar)
	VisitorCVar map[string]string

	// when this is set to true, no information will be sent to piwik
	Ignore bool

	params url.Values
}

// New implements the Matomo Tracking WEB API
//
//
func New(options ...Options) (*Matomo, error) {

	opts, err := prepareOptions(options)
	if err != nil {
		return nil, err
	}
	matomo := Matomo{
		Options: opts,
	}

	return &matomo, nil
}

func prepareOptions(options []Options) (*Options, error) {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}
	if opt.PiwikURL == "" {
		return nil, errors.New("Matomo api url must not be emty or null")
	}
	opt.PiwikURL = fixPiwikBaseURL(opt.PiwikURL)

	return &opt, nil
}

func (m *Matomo) Request(r *http.Request) {

	headers := r.Header
	if !m.Options.IgnoreDoNotTrack && headers.Get("DNT") == "1" {
		// stop execution
		return
	}

	params := make(url.Values)
	params.Set("isside", m.Options.WebsiteID)
	params.Set("rec", "1")

	proto := headers.Get("X-Forwarded-Proto")
	if proto == "" {
		if r.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}
	host := headers.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}

	params.Set("url", proto+"://"+host+r.URL.String())
	params.Set("apiv", "1")
	params.Set("urlref", headers.Get("Referer"))
	params.Set("ua", headers.Get("User-Agent"))
	params.Set("lang", headers.Get("Accept-Language"))

	ip := r.RemoteAddr
	if strings.Contains(ip, ",") {
		ipv6 := strings.Split(ip, ",")
		ip = strings.TrimPrefix(strings.TrimSpace(ipv6[0]), "::ffff:")
	}
	params.Set("token_auth", m.Options.Token)
	params.Set("cip", ip)

	p := TrackingParams{
		ActionCVar:  make(map[string]string),
		VisitorCVar: make(map[string]string),
		params:      params,
	}

	fmt.Print(p)

	// sending data to matomo
	go func() {
		return
	}()

}

func fixPiwikBaseURL(url string) string {
	if !strings.Contains(url, "/piwik.php") && !strings.Contains(url, "/proxy-piwik.php") {
		url += "/piwik.php"
	}
	return url
}
