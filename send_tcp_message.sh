#!/bin/bash

count=${1:-1}

(for i in $(seq 1 ${count}); do echo -n -e "{ \"version\": \"1.1\", \"host\": \"example.org\", \"short_message\": \"A short message\", \"level\": 5, \"_some_info\": \"foo\", \"_index\": ${i} }\0"; done) | nc -w1 127.0.0.1 12201
