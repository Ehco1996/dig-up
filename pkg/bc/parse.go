package bc

import (
	"net/http"
	"strings"
)

// b站通过这几个 cookie 来做鉴权
const (
	cookiePrefix = "-H $'cookie:"

	KEY_bili_jct   = "bili_jct"
	KEY_SESSDATA   = "SESSDATA"
	KEY_buvid3     = "buvid3"
	KEY_DedeUserID = "DedeUserID"
)

func parseCurl2Credential(curl string) []*http.Cookie {
	args := strings.Split(strings.TrimSpace(curl), "-H")
	rawCookies := ""
	for _, arg := range args {
		if strings.Contains(arg, "cookie:") {
			rawCookies = strings.TrimLeft(arg, cookiePrefix) //nolint
		}
	}
	cookies := []*http.Cookie{}

	header := http.Header{}
	header.Add("Cookie", rawCookies)
	request := http.Request{Header: header}
	for _, c := range request.Cookies() {
		if c.Name == KEY_bili_jct || c.Name == KEY_SESSDATA || c.Name == KEY_buvid3 || c.Name == KEY_DedeUserID {
			cookies = append(cookies, c)
		}
	}
	return cookies
}
