package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"spproxy/internal/configs"
	"strings"
	"time"
)

type ProxyCache struct {
	CurrentPort    string
	CurrentAppName string
	ProxyMap       map[string]*httputil.ReverseProxy
}

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy
}

func NewProxyCache() *ProxyCache {
	return &ProxyCache{
		CurrentPort: "",
		ProxyMap:    make(map[string]*httputil.ReverseProxy),
	}
}

func ProxyRequestHandler(url *url.URL, resource configs.Resource, proxyCache *ProxyCache) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		originalURL := r.URL.String()

		// Update the headers to allow for SSL redirection
		if resource.Endpoint == "/" && resource.Port != proxyCache.CurrentPort {
			r.URL.Host = strings.Replace(url.Host, resource.Port, proxyCache.CurrentPort, 1)
			r.Host = strings.Replace(url.Host, resource.Port, proxyCache.CurrentPort, 1)
		} else {
			r.URL.Host = url.Host
			r.Host = url.Host
		}
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

		// Trim reverseProxyRoutePrefix
		path := r.URL.Path
		trimLen := len(resource.Endpoint)
		if path[:trimLen] == resource.Endpoint {
			r.URL.Path = path[trimLen:]
		}

		// Update the sticky proxy port
		if resource.Endpoint != "/" || proxyCache.CurrentPort == "" {
			proxyCache.CurrentPort = resource.Port
			proxyCache.CurrentAppName = resource.Name

			createProxy(url, proxyCache)

			redirectPath := originalURL[len(resource.Endpoint)-1:]
			fmt.Printf("[%s] *** Browser redirect for %s to %s\n", timeString(), resource.Name, redirectPath)
			http.Redirect(w, r, redirectPath, http.StatusSeeOther)
			return
		}

		createProxy(url, proxyCache)

		// Send the request to the proxy for the sticky port
		proxyCache.ProxyMap[proxyCache.CurrentPort].ServeHTTP(w, r)
		fmt.Printf("[%s] %s Request %s %s ==> %s\n", timeString(), proxyCache.CurrentAppName, r.Method, originalURL, r.URL)
	}
}

func createProxy(url *url.URL, proxyCache *ProxyCache) {
	if proxyCache.ProxyMap[proxyCache.CurrentPort] == nil {
		fmt.Printf("[%s] *** Create new proxy for %s on %v ***\n", timeString(), proxyCache.CurrentAppName, url)
		proxyCache.ProxyMap[proxyCache.CurrentPort] = NewProxy(url)
	}
}

func timeString() string {
	loc, err := time.LoadLocation(time.Local.String())
	if err != nil {
		log.Fatal("Failed to load timezone location")
	}
	return time.Now().In(loc).Format("15:04:05")
}
