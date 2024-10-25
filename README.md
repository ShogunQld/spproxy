# spproxy - Sticky Port Proxy

This is a dev tool that lets you run multiple web apps on the / url path and jump between them seamlessly without running into CORS issues.  It is a reverse proxy that will update the port that the / route targets each time you hit a different route with a port attached to it and redirects the browser back to the / route.

It is NOT:
- For production use
- Useful for anything other than running things on localhost for development purposes
- For use cases that require more than one session at a time (the / route changing targets will break the other session)
