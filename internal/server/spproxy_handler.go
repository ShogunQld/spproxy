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
	CurrentDestURL string
	ProxyMap       map[string]*ProxyRoute
}

func (pc *ProxyCache) CurrentPort() string {
	if pc.ProxyMap[pc.CurrentDestURL] == nil {
		return ""
	}
	return pc.ProxyMap[pc.CurrentDestURL].Route.Port
}

func (pc *ProxyCache) CurrentAppName() string {
	if pc.ProxyMap[pc.CurrentDestURL] == nil {
		return ""
	}
	return pc.ProxyMap[pc.CurrentDestURL].Route.Name
}

func (pc *ProxyCache) CurrentEndpoint() string {
	if pc.ProxyMap[pc.CurrentDestURL] == nil {
		return ""
	}
	return pc.ProxyMap[pc.CurrentDestURL].Route.Endpoint
}

type ProxyRoute struct {
	ReverseProxy *httputil.ReverseProxy
	Route        *configs.Route
}

func NewProxy(target *url.URL, route *configs.Route) *ProxyRoute {
	proxy := httputil.NewSingleHostReverseProxy(target)
	return &ProxyRoute{
		ReverseProxy: proxy,
		Route:        route,
	}
}

func NewProxyCache() *ProxyCache {
	return &ProxyCache{
		CurrentDestURL: "",
		ProxyMap:       make(map[string]*ProxyRoute),
	}
}

func ProxyRequestHandler(route *configs.Route, proxyCache *ProxyCache) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse(route.Destination_URL)
		if err != nil {
			log.Fatalf("Failed to create ProxyRequestHandler for %s : %v", route.Name, err)
		}

		originalURL := r.URL.String()

		// Update the headers to allow for SSL redirection
		if route.Endpoint == "/" && route.Port != proxyCache.CurrentPort() {
			r.URL.Host = strings.Replace(url.Host, route.Port, proxyCache.CurrentPort(), 1)
			r.Host = strings.Replace(url.Host, route.Port, proxyCache.CurrentPort(), 1)
		} else {
			r.URL.Host = url.Host
			r.Host = url.Host
		}
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

		// Trim reverseProxyRoutePrefix
		path := r.URL.Path
		trimLen := len(route.Endpoint)
		if path[:trimLen] == route.Endpoint {
			r.URL.Path = path[trimLen:]
		}

		// Update the sticky proxy port
		if (route.Endpoint != "/" || proxyCache.CurrentPort() == "") && route.Port != "" {
			proxyCache.CurrentDestURL = route.Destination_URL
		}

		createProxy(url, proxyCache, route)

		// If a sticky port is defined for the current route, remove the endpoint url segment and
		// redirect the client to the sticy port route via a HTTP 303 Redirect
		if route.Endpoint != "/" && route.Port != "" {
			redirectPath := originalURL[len(route.Endpoint)-1:]
			fmt.Printf("[%s] *** Browser redirect for %s to %s\n", timeString(), route.Name, redirectPath)
			http.Redirect(w, r, redirectPath, http.StatusSeeOther)
			return
		}

		destUrl := proxyCache.CurrentDestURL
		appName := proxyCache.CurrentAppName()
		if route.Endpoint != "/" && route.Port == "" {
			destUrl = route.Destination_URL
			appName = route.Name
		}

		// Send the request to the proxy for the sticky port
		proxyCache.ProxyMap[destUrl].ReverseProxy.ServeHTTP(w, r)
		fmt.Printf("[%s] %s Request %s %s ==> %s\n", timeString(), appName, r.Method, originalURL, r.URL)
	}
}

func createProxy(url *url.URL, proxyCache *ProxyCache, route *configs.Route) {
	if proxyCache.ProxyMap[route.Destination_URL] == nil {
		fmt.Printf("[%s] *** Create new proxy for %s on %v ***\n", timeString(), route.Name, url)
		proxyCache.ProxyMap[route.Destination_URL] = NewProxy(url, route)
	}
}

func timeString() string {
	loc, err := time.LoadLocation(time.Local.String())
	if err != nil {
		log.Fatal("Failed to load timezone location")
	}
	return time.Now().In(loc).Format("15:04:05")
}
