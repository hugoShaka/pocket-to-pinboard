## Pocket to Pinboard migration script

[![No Maintenance Intended](http://unmaintained.tech/badge.svg)](http://unmaintained.tech/)

This repo contains a small program that logs into Pocket, retrieves every saved entry, and posts them to Pinboard.

This tries to retain titles, tags, descriptions and save date. The "favourite" attribute is not kept because it's 
not available via Pinboard's v1 API, there's no Pinboard v2 API off the shelf, and I'm too lazy to write one, I just 
want to migrate my saved content.

### Usage

You can find the pinboard api key here: https://pinboard.in/settings/password

You can get your pocket app consumer key here: https://getpocket.com/developer/apps/

Or create a new one here: https://getpocket.com/developer/apps/new

```shell
go run ./cmd/pocket-to-pinboard/ migrate --pinboard-api-key "<pinboardName>:<random part>" --pocket-consumer-key "<key>"
```
