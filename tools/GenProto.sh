#!/usr/bin/env bash
protoc --proto_path=./Proto/src --go_out=../msg ./Proto/src/*.proto 
# protoc --proto_path=./Proto/src --csharp_out=D:/FF/Test/TestClient/Assets/Scripts/Net/proto ./Proto/src/*.proto

python3 gen_proto.py

read -p "Press any key to continue."