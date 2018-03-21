#!/usr/bin/env bash -v
echo 'hello from bash'
echo 'running wget -O /tmp/wgetout http://google.com'
wget -O /tmp/wgetout http://google.com/ 2>&1
echo 'running `cat ~/.aws/credentials` from inside docker'
cat ~/.aws/credentials || echo "No ~/.aws/credentials file found"
echo 'bye from bash'
