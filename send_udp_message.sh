#!/bin/bash

echo -n '{ "version": "1.1", "host": "example.org", "short_message": "A short message", "level": 5, "_some_info": "foo" }' | nc -w1 -u 127.0.0.1 12201
