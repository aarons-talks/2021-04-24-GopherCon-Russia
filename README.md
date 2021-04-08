# GopherCon Russia 2021

This repository contains slides and code for my talk at [GopherCon Russia 2021](https://www.gophercon-russia.ru/en).

## Slides

The [`slides/`](./slides) directory has the slides for the talk. They are in [reveal JS](https://revealjs.com/) format. In order to run them, you need a static web server. Support is provided out of the box for [Caddy](https://caddyserver.com/). To use it, make sure you have Caddy v2 and run the following from the root of this repository (not the slides directory):

```shell
caddy run
```

## Code

The other 3 directories are Go code:

- [`origin/`](./origin) - a simple HTTP server intended to be the "application" that we're proxying
- [`scaler/`](./scaler) - a simple HTTP server intended to represent the system that scales replicas of the origin up and down. It also reports the number of replicas of the origin.
    >The scaler doesn't actually scale replicas in this demo. It just reports a fake number of replicas
- [`proxy/`](./proxy) - the service that accepts HTTP requests and forwards them to the origin. The proxy uses the scaler to determine how many replicas of the origin exist (it holds the request if there are no replicas) and intelligently forwards requests to the origin based on whether it can establish a TCP connection with it.
