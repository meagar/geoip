# geoip

A dead-simple CLI for geo-locating up IP addresses using MaxMindb's GeoLite2-City.

## Install

Use `go install`:

```
$ go install github.com/meagar/geoip@latest
```

## Usage:

Use `geoip <ip address>`

```
$ geoip 8.8.8.8
United States

$ geoip 194.60.38.225
Kensington, England, Royal Kensington and Chelsea, United Kingdom
...

