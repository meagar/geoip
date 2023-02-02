# geoip

A thin CLI wrapper around MaxMindb's GeoLite2-City database.

## Install

Use `go install`:

```
$ go install github.com/meagar/geoip@latest
```

## Usage:

Use `geoip [-locale=LOCALE] [IP_ADDRESS]`

```
$ geoip 8.8.8.8
United States

$ geoip 194.60.38.225
Kensington, England, Royal Kensington and Chelsea, United Kingdom

$ geoip --locale=ru 194.60.38.225
Кенсингтон, Англия, Великобритания
```

