package server

import (
	"fmt"
	"net/http"
	"spproxy/internal/configs"
	"strings"
)

func StatusRequestHandler(config *configs.Configuration, proxyCache *ProxyCache) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("**** Status Endpoint ****")
		var sb strings.Builder
		sb.WriteString(`<html>
		<head>
		<title>Sticky Port Proxy Status</title>
		<style>
			table {
				font-family: Arial, Helvetica, sans-serif;
				border-collapse: collapse;
			}

			td, th {
				border: 1px solid #ddd;
				padding: 8px;
			}

			tr:nth-child(even){background-color: #f2f2f2;}

			th {
				padding-top: 12px;
				padding-bottom: 12px;
				text-align: left;
				background-color: #555555;
				color: white;
			}
		</style>
		</head>
		</body>
		<center>`)
		sb.WriteString("<h1>Sticky Port Proxy Status</h1>\n")
		if proxyCache.ProxyMap[proxyCache.CurrentDestURL] == nil {
			sb.WriteString("<strong>No Sticky Port Set</strong>\n")
		} else {
			sb.WriteString(fmt.Sprintf("Current Sticky Port <strong>%s</strong> for <strong>%s</strong> redirects <strong>%s</strong> to <strong>%s</strong><p>\n",
				proxyCache.CurrentPort(), proxyCache.CurrentAppName(), proxyCache.CurrentEndpoint(), proxyCache.CurrentDestURL))
		}
		sb.WriteString("<h2>Routes</h2>\nClick on the URL link to update the current stick port to that route.<p>\n")
		sb.WriteString("<table>\n<tr><th>Route Name</th><th>Sticky Port</th><th>Endpoint</th><th>URL</th></tr>\n")
		for _, r := range config.Routes {
			sb.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
				r.Name, r.Port, r.Endpoint, buildUrlLink(config, r)))
		}
		sb.WriteString("</table>\n</center>\n</body>\n")
		w.Write([]byte(sb.String()))
	}
}

func buildUrlLink(config *configs.Configuration, route configs.Route) string {
	if route.Endpoint == "/" {
		return route.Destination_URL
	}
	url := route.Endpoint
	if url[0] == '/' {
		url = url[1:]
	}
	if len(url) > 0 && (url[len(url)-1] != '/') {
		url = url + "/"
	}
	return fmt.Sprintf("<a href=\"http://%s:%s/%sstatus\">%s</a>",
		config.Server.Host, config.Server.Listen_port, url, route.Destination_URL)
}
