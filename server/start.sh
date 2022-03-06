#!/bin/bash
## 退出时清理资源，先删除子进程
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8001 &
./server -port=8002 &
./server -port=8003

sleep 2

## 防止 bash 进程退出
wait
