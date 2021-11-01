#!/bin/sh
#protoc --go_out=plugins=grpc,paths=source_relative:. pkg/message.proto
protoc --go_out=plugins=grpc,paths=source_relative:. $@
