#!/bin/bash

##To generate protobuf files to .go file only
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/addressbook.proto

##To generate grpc service file only
protoc --go-grpc_out=. greet/greetpb/greet.proto

##To generate both grpc service and protobuf message files///run from project root
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative greet/greetpb/greet.proto

## For calculator service
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative calculator/calculatorpb/calculator.proto
