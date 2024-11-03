
# spproxy - Sticky Port Proxy

This is a dev tool that lets you run multiple web apps on the `/` route and jump between them seamlessly without running into CORS issues.  It is a reverse proxy that will update the port that the `/` route targets each time you hit a different route with a port attached to it and redirects the browser back to the `/` route.

It is NOT:
- For production use
- Useful for anything other than running things on localhost for development purposes
- For use cases that require more than one session at a time (the `/` route changing targets will break the other session)

## How it works

Assume you are running spproxy on `localhost:8080`

If you define a sticky route of `/app1/` that points to `localhost:9999`, requests to `localhost:8080/app1/url/path` would respond with a HTTP 303 redirect to `localhost:8080/url/path`.  When you follow `localhost:8080/url/path` the stick port proxy will forward your request to `localhost:9999/url/path` and pass back the response.  Then future requests to `localhost:8080` that don't match any other defined routes will also be forwarded to `localhost:9999`.

If you define a non-sticky route `/nostick/` that points to `localhost:1111`, requests to `localhost:8080/nostick/` will be forwarded to `localhost:1111`.  The sticky port WON'T be updated, so future requests to `localhost:8080` that don't match any other defined routes will NOT be forwarded to `localhost:1111`.  Instead they will go to the last sticky port route (i.e. `localhost:9999`)

**Example:**

With a config that defines the following:
- /app1/ -> localhost:1111 (with a sticky port of 1111)
- /app2/ -> localhost:2222 (with a sticky port of 2222)
- /app3/ -> localhost:3333 (no sticky port)

From a web browser that follows the redirects, requests will resolve as follows:
| Request Path | Redirect To |
|--|--|
| localhost:8080/app1/test | localhost:1111/test |
| localhost:8080/data | localhost:1111/data |
| localhost:8080/app2/test | localhost:2222/test |
| localhost:8080/data | localhost:2222/data |
| localhost:8080/app3/test | localhost:3333/test |
| localhost:8080/data | localhost:2222/data |

Note: The final request to `localhost:8080/test` still hits `localhost:2222/test` because the `/app3/` ports isn't sticky, so the base sticky route still points to app2.

## Configuration

By default spproxy loads it configuration from `config.json` in the same directory as the executable.  However you can pass the path to your configuration file as the first command line parameter to spproxy.

The following example configuration does the following:
- The `server` section defines the server host name and port that ssproxy will run on.
- Sticky: The sticky port that will default to `localhost.8081` when spproxy starts.
- Server1: Redirect requests on the `/server1/` route to `localhost:8081` and update the Sticky route to target port `8081`.
- Server2: Redirect requests on the `/server2/` route to `localhost:8082` and update the Sticky route to target port `8082`.
- Server3: Redirect requests on the `/server3/` route to `localhost:8083` but don't update Sticky route.

```
{
  "server": {
    "host": "localhost",
    "listen_port": "8080"
  },
  "routes": [
    {
      "name": "Sticky",
      "endpoint": "/",
      "destination_url": "http://localhost:8081",
    },
    {
      "name": "Server1",
      "endpoint": "/server1/",
      "destination_url": "http://localhost:8081",
      "port": "8081"
    },
    {
      "name": "Server2",
      "endpoint": "/server2/",
      "destination_url": "http://localhost:8082",
      "port": "8082"
    },
    {
      "name": "Server3",
      "endpoint": "/server3/",
      "destination_url": "http://localhost:8083"
    }
  ]
}
```

## Status Page

When running you can visit the `/status` route and see the current sticky port as well as all the routes.  This page also allows up to update the current sticky page.

## Limitations

The Sticky Port Proxy has the following limitations:
- It is designed as a dev tool and as such is not to be used in a production environment.
- It only supports a single session at a time, as it doesn't distinguish between were requests come from.
- It will only work with web apps that expect to service requests on the `/` route, but redirect to each other using different base routes.
- This is was learnt to use go, so it probably does lots of things in non-optimal ways that aren't done in the go idiomatic way.

## Why build this type of proxy

Two reasons:
- My work had multiple web apps that fit the third limitation above that required constantly stopping and starting apps to jump between them when running them locally, which was painful.
- More importantly, it was a good excuse to learn go.

## Air hot reloading

This project contains an air configuration to allow hot reloading during development.

To install Air, run `make install`

To run spproxy using Air, run `air <config path>`
