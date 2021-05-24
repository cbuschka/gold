#!/bin/bash

first=${1}
limit=${2}

curl -v --unix-socket ./tmp/golfd.sock -X GET "http:/host/messages?begin=${1}&limit=${2}"
