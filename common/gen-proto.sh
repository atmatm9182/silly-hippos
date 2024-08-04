#!/bin/env sh

set -xe

protoc -I=./proto/ --go_out=. --go_opt=paths=source_relative $(find ./proto/ -iname "*.proto")
