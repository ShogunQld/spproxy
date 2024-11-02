package server

import (
	"fmt"
	"net/http"
	"spproxy/internal/configs"
)

// Run starts server and listens on defined port
func Run(config *configs.Configuration) error {
	// Create a new router
	mux := http.NewServeMux()

	proxyCache := NewProxyCache()

	mux.HandleFunc("/status", StatusRequestHandler(config, proxyCache))

	// Iterate through the configuration and register the routes into the proxy cache
	for _, route := range config.Routes {
		mux.HandleFunc(route.Endpoint, ProxyRequestHandler(&route, proxyCache))
		if route.Endpoint == "/" {
			proxyCache.CurrentDestURL = route.Endpoint
		}
	}

	// Run the proxy server
	fmt.Printf("Start Proxy on %s:%s\n", config.Server.Host, config.Server.Listen_port)
	if err := http.ListenAndServe(config.Server.Host+":"+config.Server.Listen_port, mux); err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}

	return nil
}
