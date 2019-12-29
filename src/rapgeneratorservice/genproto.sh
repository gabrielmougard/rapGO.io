#!/bin/bash -e

python -m grpc_tools.protoc -I../../protobuf --python_out=./genproto --grpc_python_out=./genproto ../../protobuf/rapgo.proto
python -m grpc_tools.protoc -I../../protobuf/grpc/health/v1 --python_out=./genproto --grpc_python_out=./genproto ../../protobuf/grpc/health/v1/health.proto