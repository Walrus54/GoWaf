package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"waf/internal/waf"
)

func StartServer(waf *waf.WAF) error {
	target, err := url.Parse(waf.Config.TargetURL)
	if err != nil {
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !isRequestAllowed(waf, r) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("Request blocked by WAF"))
			return
		}

		proxy.ServeHTTP(w, r)
	})

	return http.ListenAndServe(waf.Config.WafPort, nil)
}

func isRequestAllowed(waf *waf.WAF, r *http.Request) bool {
	// Проверка URL и GET-параметров
	if !checkString(r.URL.String(), waf, r) {
		return false
	}

	// Проверка тела запроса для POST/PUT
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		err := r.ParseForm()
		if err != nil {
			waf.Logger.LogError("Failed to parse form: " + err.Error())
			return false
		}

		// Проверка POST-параметров
		if !checkString(r.Form.Encode(), waf, r) {
			return false
		}
	}

	// Проверка заголовков
	for _, h := range r.Header {
		if !checkString(strings.Join(h, ","), waf, r) {
			return false
		}
	}

	return true
}

func checkString(s string, waf *waf.WAF, r *http.Request) bool {
	for _, rule := range waf.Rules {
		if rule.Pattern.MatchString(s) {
			waf.Logger.LogBlockedRequest(r, rule.Name)
			return false
		}
	}
	return true
}
