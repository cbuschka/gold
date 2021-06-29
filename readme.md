# gold - gelf log daemon written in go
[![Build](https://github.com/cbuschka/gold/workflows/build/badge.svg)](https://github.com/cbuschka/gold) [![License](https://img.shields.io/github/license/cbuschka/gold.svg)](https://github.com/cbuschka/gold/blob/main/license.txt)

### WIP!

## Features

* receives gelf messages on udp, tcp and http
* stores log messages in ~~boltdb~~ ~~badger~~ ~~leveldb~~ [pebble by cockroach labs](https://github.com/cockroachdb/pebble)
* adds generated uuid as \_id attribute to gelf message
* adds \_received\_timestamp and \_sender\_host attributes to gelf message
* exports rest api via unix domain socket for querying and control

## Planned Features

* query client that uses the rest api
* archive and expunge old log data

## Configuration

gold.conf.json

```json
{
  "dataDir": "./data/db",
  "commandSocketPath": "./run/gold.sock",
  "gelfUdpListeners": [
    "127.0.0.1:12201"
  ],
  "gelfTcpListeners": [
    "127.0.0.1:12201"
  ],
  "gelfHttpListeners": [
    "127.0.0.1:8080"
  ]
}
```

## References

* [GELF doc](https://docs.graylog.org/en/4.0/pages/gelf.html)

## License
Copyright (c) 2020-2021 by [Cornelius Buschka](https://github.com/cbuschka).

[Apache License, Version 2.0](./license.txt)
