package httpcrud

import (
	"net/http"
	"regexp"
)

// CookieValues can be used to read the values from the Cookie header of the
// incoming HTTP request and indirectly set them to the strings pointed to by
// the map's values, using the map's keys to match against the Cookies' names.
type CookieValues map[string]*string

// ReadHeader implements the HeaderReader interface.
func (rr CookieValues) ReadHeader(header http.Header) error {
	cc := (&http.Request{Header: header}).Cookies()
	if len(cc) == 0 {
		return nil
	}

	for key, val := range rr {
		for i := 0; i < len(cc); i++ {
			if cc[i].Name == key {
				*val = cc[i].Value
				break
			}
		}
	}
	return nil
}

// IPAddress can be used to read the IP address value from the header of the
// incoming HTTP request and set it to the string pointed to by Val.
type IPAddress struct {
	// Pointer to the string which should be set to the ip address.
	Val *string
}

// ReadHeader implements the HeaderReader interface.
func (rr IPAddress) ReadHeader(header http.Header) error {
	if ip := header.Get("X-Forwarded-For"); len(ip) > 0 {
		*rr.Val = ip
	} else if ip := header.Get("X-Real-Ip"); len(ip) > 0 {
		*rr.Val = ip
	}
	return nil
}

// UserAgent can be used to read the value from the User-Agent header of the
// incoming HTTP request and set it to the string pointed to by Val.
type UserAgent struct {
	// Pointer to the string which should be set to the user agent.
	Val *string
}

// ReadHeader implements the HeaderReader interface.
func (rr UserAgent) ReadHeader(header http.Header) error {
	if ua := header.Get("User-Agent"); len(ua) > 0 {
		*rr.Val = ua
	}
	return nil
}

var rxBearer = regexp.MustCompile(`(?i:bearer\s+)([0-9A-Za-z\-_]+)`)

// BearerToken can be used to read the bearer token value from the Authorization
// header of the incoming HTTP request and set it to the string pointed to by Val.
type BearerToken struct {
	// Pointer to the string which should be set to the bearer token.
	Val *string
}

// ReadHeader implements the HeaderReader interface.
func (rr BearerToken) ReadHeader(header http.Header) error {
	auth := header.Get("Authorization")
	if match := rxBearer.FindStringSubmatch(auth); len(match) > 1 {
		*rr.Val = match[1]
	}
	return nil
}

// SetCookie can be used to set the Set-Cookie header of the outgoing response.
type SetCookie struct {
	// The cookie to be set.
	Val *http.Cookie
}

// WriteHeader implements the HeaderWriter interface.
func (rr SetCookie) WriteHeader(header http.Header) {
	if v := rr.Val.String(); v != "" {
		header.Add("Set-Cookie", v)
	}
}
