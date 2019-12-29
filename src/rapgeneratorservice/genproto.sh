#!/bin/bash -e

python -m grpc_tools.protoc -I../../protobuf --python_out=. --grpc_python_out=. ../../protobuf/rapgo.proto