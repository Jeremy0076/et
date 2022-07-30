#!/bin/bash

echo "$(TZ=UTC-8 date +%Y-%m-%d" "%H:%M:%S) nginx conf starting"
# 拷贝nginx.conf
cp nginx.conf /opt/homebrew/etc/nginx/
# 拷贝前端页面到nginx下
cd ../frontend
cp -r /img /opt/homebrew/var/www/
cp index.html /opt/homebrew/var/www/
cp login.html /opt/homebrew/var/www/

echo "$(TZ=UTC-8 date +%Y-%m-%d" "%H:%M:%S) tcp server starting"
cd ../tcpserver
go build main.go -o tcpsvr
# 后台运行
./tcpsvr ../conf/conf.ini >&1 &

echo "$(TZ=UTC-8 date +%Y-%m-%d" "%H:%M:%S) http server starting"
cd ../httpserver
go build main.go -o httpsvr
./httpsvr ../conf/conf.ini >&1 &

