#!/bin/bash

if [ -z "${1}${2}" ]; then
  curl -v --unix-socket ./run/gold.sock -X GET "http:/host/messages"
else
  curl -v --unix-socket ./run/gold.sock -X GET "http:/host/messages?begin=${1}&limit=${2}"
fi
