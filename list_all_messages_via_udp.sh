#!/bin/bash

curl -v --unix-socket ./tmp/golfd.sock -X GET http:/host/messages
