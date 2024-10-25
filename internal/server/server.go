package server

import (
	"fmt"
	"net/http"
	"net/url"
	"spproxy/internal/configs"
)

// Run starts server and listens on defined port
func Run(config *configs.Configuration) error {
	// Create a new router
	mux := http.NewServeMux()

	// Iterate through the configuration and register the routes into the proxy cache
	proxyCache := NewProxyCache()
	for _, resource := range config.Resources {
		url, _ := url.Parse(resource.Destination_URL)
		mux.HandleFunc(resource.Endpoint, ProxyRequestHandler(url, resource, proxyCache))
	}

	// Run the proxy server
	fmt.Printf("Start Proxy on %s:%s\n", config.Server.Host, config.Server.Listen_port)
	if err := http.ListenAndServe(config.Server.Host+":"+config.Server.Listen_port, mux); err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}

	return nil
}
