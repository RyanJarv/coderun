#!/usr/bin/env bash
echo 'hello from bash'
echo 'running curl http://google.com'
wget -O /tmp/wgetout http://google.com/
#echo 'running `cat ~/.aws/credentials` from inside docker'
#cat ~/.aws/credentials || echo "No ~/.aws/credentials file found"
apt-get update
apt-get install iproute2 -y
ip addr
echo 'bye from bash'
