package hxio

import (
	"net/http"
	"regexp"
)

type CookieValue map[string]*string

func (rr CookieValue) ReadHeader(header http.Header) error {
	cc := (&http.Request{Header: header}).Cookies()
	if len(cc) == 0 {
		return nil
	}

	for k, v := range rr {
		for i := 0; i < len(cc); i++ {
			if cc[i].Name == k {
				*v = cc[i].Value
				break
			}
		}
	}
	return nil
}

type IPAddress struct {
	Val *string
}

func (rr IPAddress) ReadHeader(header http.Header) error {
	if ip := header.Get("X-Forwarded-For"); len(ip) > 0 {
		*rr.Val = ip
	} else if ip := header.Get("X-Real-Ip"); len(ip) > 0 {
		*rr.Val = ip
	}
	return nil
}

type UserAgent struct {
	Val *string
}

func (rr UserAgent) ReadHeader(header http.Header) error {
	if ua := header.Get("User-Agent"); len(ua) > 0 {
		*rr.Val = ua
	}
	return nil
}

var rxBearer = regexp.MustCompile(`(?i:bearer\s+)([0-9A-Za-z\-_]+)`)

type BearerToken struct {
	Val *string
}

func (rr BearerToken) ReadHeader(header http.Header) error {
	auth := header.Get("Authorization")
	if match := rxBearer.FindStringSubmatch(auth); len(match) > 1 {
		*rr.Val = match[1]
	}
	return nil
}

type Cookie struct {
	Val *http.Cookie
}

func (rr Cookie) WriteHeader(header http.Header) {
	if v := rr.Val.String(); v != "" {
		header.Add("Set-Cookie", v)
	}
}
