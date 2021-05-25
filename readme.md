# golf - gelf daemon written in go
[![Build](https://github.com/cbuschka/golf/workflows/build/badge.svg)](https://github.com/cbuschka/golf) [![License](https://img.shields.io/github/license/cbuschka/golf.svg)](https://github.com/cbuschka/golf/blob/main/license.txt)

### WIP!

## Features

* receives gelf messages on udp, tcp and http
* stores log messages in ~~boltdb~~ ~~badger~~ [goleveldb](https://github.com/syndtr/goleveldb)
* adds generated uuid as \_id attribute to gelf message
* adds \_received\_timestamp and \_sender\_host attributes to gelf message
* exports rest api via unix domain socket

## Planned Features

* query client that uses the rest api
* periodic dump to gzipped jsonlines file

## Configuration

golfd.conf.json

```json
{
  "dataDir": "./data/db",
  "commandSocketPath": "./run/golfd.sock",
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
