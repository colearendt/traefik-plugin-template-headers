# Traefik Plugin: Auth Delay

A [traefik](https://traefik.io/) plugin to build additional request headers based on go templates.

This is meant to help replace some of the behavior that is common with request modification middlewares in
proxies like Apache/`httpd` or `nginx`.

Data passed along to the template from the request object (`req`) for an example
URL (`https://localhost:80/some/path?query=true`):

| key                 | value                               | Example               |
|---------------------|-------------------------------------|-----------------------|
| Path                | req.URL.EscapedPath()               | /some/path            |
| Scheme              | req.URL.Scheme                      | https                 |
| Host                | req.URL.Host                        | localhost:80          |
| Method              | req.Method                          | GET                   |
| Proto               | req.Proto                           | HTTP/1.1              |
| Query               | req.URL.RawQuery                    | query=true            |
| RequestURI          | req.URL.RequestURI()                | /some/path?query=true |
| HttpXForwardedProto | req.Header.Get("X-Forwarded-Proto") | https                 |
| HttpXForwardedHost  | req.Header.Get("X-Forwarded-Host")  | localhost:80          |
| HttpHost            | req.Header.Get("Host")              | localhost:80          |



## Example Configuration

TODO

## What is a Traefik Plugin

TL;DR; A Traefik plugin is a custom middleware for Traefik.

[More on Traefik plugins is written here](https://doc.traefik.io/traefik/plugins/).

I also wrote [an init container](https://github.com/colearendt/traefik-plugin-init) that simplifies using "local"
plugins (i.e. plugins without Traefik Pilot) inside of Kubernetes.

## TODO

- Decide about semantics for modifying, overwriting, appending, or removing headers

## Thanks

Inspired by and much boilerplate
from [traefik-plugin-rewrite-headers](https://github.com/XciD/traefik-plugin-rewrite-headers), which is a fantastically
useful Traefik Plugin.
